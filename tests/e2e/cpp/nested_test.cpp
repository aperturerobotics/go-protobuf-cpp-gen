// Nested messages test - verifies nested message handling
#include <iostream>
#include <string>

#include "nested.pb.h"

#define ASSERT(cond, msg) \
    if (!(cond)) { \
        std::cerr << "FAILED: " << msg << " at " << __FILE__ << ":" << __LINE__ << std::endl; \
        return 1; \
    }

#define ASSERT_EQ(a, b, msg) ASSERT((a) == (b), msg)

int main() {
    std::cout << "Testing nested messages..." << std::endl;

    // Test Outer with nested Middle and Inner
    {
        test::nested::Outer outer;

        outer.set_id("outer-001");
        ASSERT_EQ(outer.id(), "outer-001", "outer id");

        // Access nested Middle via mutable
        auto* middle = outer.mutable_middle();
        middle->set_name("middle-name");

        // Access deeply nested Inner
        auto* inner = middle->mutable_inner();
        inner->set_value("inner-value");
        inner->set_count(42);

        // Verify values
        ASSERT_EQ(outer.middle().name(), "middle-name", "middle name");
        ASSERT_EQ(outer.middle().inner().value(), "inner-value", "inner value");
        ASSERT_EQ(outer.middle().inner().count(), 42, "inner count");

        std::cout << "  Basic nested access: OK" << std::endl;

        // Test repeated nested messages
        auto* inner1 = middle->add_inners();
        inner1->set_value("inner1");
        inner1->set_count(1);

        auto* inner2 = middle->add_inners();
        inner2->set_value("inner2");
        inner2->set_count(2);

        ASSERT_EQ(middle->inners_size(), 2, "inners size");
        ASSERT_EQ(middle->inners(0).value(), "inner1", "inners[0] value");
        ASSERT_EQ(middle->inners(1).value(), "inner2", "inners[1] value");

        std::cout << "  Repeated nested messages: OK" << std::endl;

        // Test repeated Middle
        auto* m1 = outer.add_middles();
        m1->set_name("middle1");
        auto* m2 = outer.add_middles();
        m2->set_name("middle2");

        ASSERT_EQ(outer.middles_size(), 2, "middles size");
        ASSERT_EQ(outer.middles(0).name(), "middle1", "middles[0] name");
        ASSERT_EQ(outer.middles(1).name(), "middle2", "middles[1] name");

        std::cout << "  Repeated middle messages: OK" << std::endl;
    }

    // Test Document with multiple nested types
    {
        test::nested::Document doc;

        // Set header
        auto* header = doc.mutable_header();
        header->set_title("Test Document");
        header->set_author("Test Author");
        header->set_created_at(1234567890);

        // Set body with sections
        auto* body = doc.mutable_body();
        auto* section1 = body->add_sections();
        section1->set_heading("Introduction");
        section1->set_content("This is the introduction.");

        auto* section2 = body->add_sections();
        section2->set_heading("Main Content");
        section2->set_content("This is the main content.");

        // Set footer
        auto* footer = doc.mutable_footer();
        footer->set_copyright("Copyright 2024");
        footer->set_page_count(10);

        // Verify
        ASSERT_EQ(doc.header().title(), "Test Document", "header title");
        ASSERT_EQ(doc.header().author(), "Test Author", "header author");
        ASSERT_EQ(doc.header().created_at(), 1234567890, "header created_at");

        ASSERT_EQ(doc.body().sections_size(), 2, "sections size");
        ASSERT_EQ(doc.body().sections(0).heading(), "Introduction", "section0 heading");
        ASSERT_EQ(doc.body().sections(1).heading(), "Main Content", "section1 heading");

        ASSERT_EQ(doc.footer().copyright(), "Copyright 2024", "footer copyright");
        ASSERT_EQ(doc.footer().page_count(), 10, "footer page_count");

        std::cout << "  Document structure: OK" << std::endl;
    }

    // Test deeply nested Level1-4
    {
        test::nested::Level1 level1;

        auto* level2 = level1.mutable_level2();
        auto* level3 = level2->mutable_level3();
        auto* level4 = level3->mutable_level4();

        level4->set_deep_value("deeply nested value");

        ASSERT_EQ(level1.level2().level3().level4().deep_value(), "deeply nested value", "deep nesting");

        std::cout << "  Deep nesting (4 levels): OK" << std::endl;
    }

    // Test standalone nested types
    {
        test::nested::Middle middle;  // All types at namespace level
        middle.set_name("standalone middle");

        test::nested::Inner inner;
        inner.set_value("standalone inner");
        inner.set_count(99);

        ASSERT_EQ(middle.name(), "standalone middle", "standalone middle");
        ASSERT_EQ(inner.value(), "standalone inner", "standalone inner value");
        ASSERT_EQ(inner.count(), 99, "standalone inner count");

        std::cout << "  Standalone nested types: OK" << std::endl;
    }

    // Test clear on nested messages
    {
        test::nested::Outer outer;
        outer.set_id("test");
        outer.mutable_middle()->set_name("middle");
        outer.mutable_middle()->mutable_inner()->set_value("inner");

        outer.clear_middle();

        // After clear, accessing should give empty/default values
        ASSERT_EQ(outer.middle().name(), "", "cleared middle name");
        ASSERT_EQ(outer.middle().inner().value(), "", "cleared inner value");

        std::cout << "  Clear nested message: OK" << std::endl;
    }

    // Test copy with nested messages
    {
        test::nested::Outer outer1;
        outer1.set_id("original");
        outer1.mutable_middle()->set_name("original-middle");
        outer1.mutable_middle()->mutable_inner()->set_value("original-inner");

        test::nested::Outer outer2(outer1);

        ASSERT_EQ(outer2.id(), "original", "copy id");
        ASSERT_EQ(outer2.middle().name(), "original-middle", "copy middle name");
        ASSERT_EQ(outer2.middle().inner().value(), "original-inner", "copy inner value");

        // Modify copy, original should be unchanged
        outer2.set_id("modified");
        outer2.mutable_middle()->set_name("modified-middle");

        ASSERT_EQ(outer1.id(), "original", "original id after copy modification");
        ASSERT_EQ(outer1.middle().name(), "original-middle", "original middle after copy modification");

        std::cout << "  Copy with nested messages: OK" << std::endl;
    }

    std::cout << "All nested message tests passed!" << std::endl;
    return 0;
}
