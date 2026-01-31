// Serialization test - verifies serialize/parse round-trip
#include <iostream>
#include <string>
#include <cmath>

#include "scalars.pb.h"
#include "enums.pb.h"
#include "repeated.pb.h"
#include "nested.pb.h"
#include "maps.pb.h"

#define ASSERT(cond, msg) \
    if (!(cond)) { \
        std::cerr << "FAILED: " << msg << " at " << __FILE__ << ":" << __LINE__ << std::endl; \
        return 1; \
    }

#define ASSERT_EQ(a, b, msg) ASSERT((a) == (b), msg)
#define ASSERT_NEAR(a, b, eps, msg) ASSERT(std::abs((a) - (b)) < (eps), msg)

// Helper to test round-trip
template<typename T>
bool roundTrip(const T& original, T& parsed) {
    std::string serialized;
    if (!original.SerializeToString(&serialized)) {
        std::cerr << "  SerializeToString failed" << std::endl;
        return false;
    }
    if (!parsed.ParseFromString(serialized)) {
        std::cerr << "  ParseFromString failed" << std::endl;
        return false;
    }
    return true;
}

int main() {
    std::cout << "Testing serialization round-trip..." << std::endl;

    // Test scalar types round-trip
    {
        test::scalars::AllScalars original;
        original.set_bool_val(true);
        original.set_int32_val(-12345);
        original.set_int64_val(-9876543210LL);
        original.set_uint32_val(12345u);
        original.set_uint64_val(9876543210ULL);
        original.set_sint32_val(-54321);
        original.set_sint64_val(-1234567890LL);
        original.set_fixed32_val(0xDEADBEEF);
        original.set_fixed64_val(0xDEADBEEFCAFEBABEULL);
        original.set_sfixed32_val(-0x12345678);
        original.set_sfixed64_val(-0x123456789ABCDEFLL);
        original.set_float_val(3.14159f);
        original.set_double_val(2.718281828);
        original.set_string_val("hello world");
        original.set_bytes_val("binary\x00data");

        test::scalars::AllScalars parsed;
        ASSERT(roundTrip(original, parsed), "AllScalars round-trip");

        ASSERT_EQ(parsed.bool_val(), true, "parsed bool");
        ASSERT_EQ(parsed.int32_val(), -12345, "parsed int32");
        ASSERT_EQ(parsed.int64_val(), -9876543210LL, "parsed int64");
        ASSERT_EQ(parsed.uint32_val(), 12345u, "parsed uint32");
        ASSERT_EQ(parsed.uint64_val(), 9876543210ULL, "parsed uint64");
        ASSERT_EQ(parsed.sint32_val(), -54321, "parsed sint32");
        ASSERT_EQ(parsed.sint64_val(), -1234567890LL, "parsed sint64");
        ASSERT_EQ(parsed.fixed32_val(), 0xDEADBEEF, "parsed fixed32");
        ASSERT_EQ(parsed.fixed64_val(), 0xDEADBEEFCAFEBABEULL, "parsed fixed64");
        ASSERT_EQ(parsed.sfixed32_val(), -0x12345678, "parsed sfixed32");
        ASSERT_EQ(parsed.sfixed64_val(), -0x123456789ABCDEFLL, "parsed sfixed64");
        ASSERT_NEAR(parsed.float_val(), 3.14159f, 0.0001f, "parsed float");
        ASSERT_NEAR(parsed.double_val(), 2.718281828, 0.0001, "parsed double");
        ASSERT_EQ(parsed.string_val(), "hello world", "parsed string");

        std::cout << "  Scalar types: OK" << std::endl;
    }

    // Test enum round-trip
    {
        test::enums::EnumMessage original;
        original.set_status(test::enums::STATUS_ACTIVE);
        original.set_priority(test::enums::PRIORITY_HIGH);
        original.add_statuses(test::enums::STATUS_PENDING);
        original.add_statuses(test::enums::STATUS_INACTIVE);

        test::enums::EnumMessage parsed;
        ASSERT(roundTrip(original, parsed), "EnumMessage round-trip");

        ASSERT_EQ(parsed.status(), test::enums::STATUS_ACTIVE, "parsed status");
        ASSERT_EQ(parsed.priority(), test::enums::PRIORITY_HIGH, "parsed priority");
        ASSERT_EQ(parsed.statuses_size(), 2, "parsed statuses size");
        ASSERT_EQ(parsed.statuses(0), test::enums::STATUS_PENDING, "parsed statuses[0]");
        ASSERT_EQ(parsed.statuses(1), test::enums::STATUS_INACTIVE, "parsed statuses[1]");

        std::cout << "  Enum types: OK" << std::endl;
    }

    // Test repeated fields round-trip
    {
        test::repeated::RepeatedScalars original;
        original.add_bool_vals(true);
        original.add_bool_vals(false);
        original.add_int32_vals(1);
        original.add_int32_vals(2);
        original.add_int32_vals(3);
        original.add_string_vals("a");
        original.add_string_vals("b");
        original.add_string_vals("c");
        original.add_double_vals(1.1);
        original.add_double_vals(2.2);

        test::repeated::RepeatedScalars parsed;
        ASSERT(roundTrip(original, parsed), "RepeatedScalars round-trip");

        ASSERT_EQ(parsed.bool_vals_size(), 2, "parsed bool_vals size");
        ASSERT_EQ(parsed.bool_vals(0), true, "parsed bool_vals[0]");
        ASSERT_EQ(parsed.bool_vals(1), false, "parsed bool_vals[1]");
        ASSERT_EQ(parsed.int32_vals_size(), 3, "parsed int32_vals size");
        ASSERT_EQ(parsed.int32_vals(0), 1, "parsed int32_vals[0]");
        ASSERT_EQ(parsed.string_vals_size(), 3, "parsed string_vals size");
        ASSERT_EQ(parsed.string_vals(0), "a", "parsed string_vals[0]");

        std::cout << "  Repeated scalars: OK" << std::endl;
    }

    // Test repeated messages round-trip
    {
        test::repeated::RepeatedMessages original;

        auto* item1 = original.add_items();
        item1->set_name("Item 1");
        item1->set_quantity(10);
        item1->set_price(9.99);

        auto* item2 = original.add_items();
        item2->set_name("Item 2");
        item2->set_quantity(5);
        item2->set_price(19.99);

        test::repeated::RepeatedMessages parsed;
        ASSERT(roundTrip(original, parsed), "RepeatedMessages round-trip");

        ASSERT_EQ(parsed.items_size(), 2, "parsed items size");
        ASSERT_EQ(parsed.items(0).name(), "Item 1", "parsed item 0 name");
        ASSERT_EQ(parsed.items(0).quantity(), 10, "parsed item 0 quantity");
        ASSERT_NEAR(parsed.items(0).price(), 9.99, 0.001, "parsed item 0 price");
        ASSERT_EQ(parsed.items(1).name(), "Item 2", "parsed item 1 name");

        std::cout << "  Repeated messages: OK" << std::endl;
    }

    // Test nested messages round-trip
    {
        test::nested::Outer original;
        original.set_id("outer-001");

        auto* middle = original.mutable_middle();
        middle->set_name("middle-name");

        auto* inner = middle->mutable_inner();
        inner->set_value("inner-value");
        inner->set_count(42);

        auto* inner1 = middle->add_inners();
        inner1->set_value("inner1");
        inner1->set_count(1);

        auto* inner2 = middle->add_inners();
        inner2->set_value("inner2");
        inner2->set_count(2);

        test::nested::Outer parsed;
        ASSERT(roundTrip(original, parsed), "Outer round-trip");

        ASSERT_EQ(parsed.id(), "outer-001", "parsed outer id");
        ASSERT_EQ(parsed.middle().name(), "middle-name", "parsed middle name");
        ASSERT_EQ(parsed.middle().inner().value(), "inner-value", "parsed inner value");
        ASSERT_EQ(parsed.middle().inner().count(), 42, "parsed inner count");
        ASSERT_EQ(parsed.middle().inners_size(), 2, "parsed inners size");
        ASSERT_EQ(parsed.middle().inners(0).value(), "inner1", "parsed inners[0] value");
        ASSERT_EQ(parsed.middle().inners(1).value(), "inner2", "parsed inners[1] value");

        std::cout << "  Nested messages: OK" << std::endl;
    }

    // Test maps round-trip
    {
        test::maps::MapScalars original;

        (*original.mutable_string_to_string())["key1"] = "value1";
        (*original.mutable_string_to_string())["key2"] = "value2";

        (*original.mutable_string_to_int32())["score1"] = 100;
        (*original.mutable_string_to_int32())["score2"] = 200;

        (*original.mutable_int32_to_string())[1] = "one";
        (*original.mutable_int32_to_string())[2] = "two";

        test::maps::MapScalars parsed;
        ASSERT(roundTrip(original, parsed), "MapScalars round-trip");

        ASSERT_EQ(parsed.string_to_string_size(), 2, "parsed string_to_string size");
        ASSERT_EQ(parsed.string_to_string().at("key1"), "value1", "parsed key1");
        ASSERT_EQ(parsed.string_to_string().at("key2"), "value2", "parsed key2");

        ASSERT_EQ(parsed.string_to_int32_size(), 2, "parsed string_to_int32 size");
        ASSERT_EQ(parsed.string_to_int32().at("score1"), 100, "parsed score1");

        ASSERT_EQ(parsed.int32_to_string().at(1), "one", "parsed 1->one");

        std::cout << "  Map scalars: OK" << std::endl;
    }

    // Test map with message values round-trip
    {
        test::maps::MapMessages original;

        test::maps::MapValue val1;
        val1.set_name("value1");
        val1.set_score(100);
        (*original.mutable_string_to_message())["key1"] = val1;

        test::maps::MapValue val2;
        val2.set_name("value2");
        val2.set_score(200);
        (*original.mutable_string_to_message())["key2"] = val2;

        test::maps::MapMessages parsed;
        ASSERT(roundTrip(original, parsed), "MapMessages round-trip");

        ASSERT_EQ(parsed.string_to_message_size(), 2, "parsed string_to_message size");
        ASSERT_EQ(parsed.string_to_message().at("key1").name(), "value1", "parsed key1 name");
        ASSERT_EQ(parsed.string_to_message().at("key1").score(), 100, "parsed key1 score");
        ASSERT_EQ(parsed.string_to_message().at("key2").name(), "value2", "parsed key2 name");

        std::cout << "  Map messages: OK" << std::endl;
    }

    // Test empty message round-trip
    {
        test::scalars::AllScalars original;
        test::scalars::AllScalars parsed;
        ASSERT(roundTrip(original, parsed), "Empty message round-trip");

        ASSERT_EQ(parsed.bool_val(), false, "empty bool");
        ASSERT_EQ(parsed.int32_val(), 0, "empty int32");
        ASSERT_EQ(parsed.string_val(), "", "empty string");

        std::cout << "  Empty message: OK" << std::endl;
    }

    // Test ByteSizeLong
    {
        test::scalars::AllScalars msg;
        msg.set_int32_val(123);
        msg.set_string_val("hello");

        size_t size = msg.ByteSizeLong();
        ASSERT(size > 0, "ByteSizeLong > 0");

        std::string serialized;
        msg.SerializeToString(&serialized);
        ASSERT_EQ(serialized.size(), size, "serialized size matches ByteSizeLong");

        std::cout << "  ByteSizeLong: OK" << std::endl;
    }

    std::cout << "All serialization tests passed!" << std::endl;
    return 0;
}
