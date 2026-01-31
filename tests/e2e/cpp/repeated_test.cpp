// Repeated fields test - verifies repeated field handling
#include <iostream>
#include <string>
#include <cmath>

#include "repeated.pb.h"

#define ASSERT(cond, msg) \
    if (!(cond)) { \
        std::cerr << "FAILED: " << msg << " at " << __FILE__ << ":" << __LINE__ << std::endl; \
        return 1; \
    }

#define ASSERT_EQ(a, b, msg) ASSERT((a) == (b), msg)
#define ASSERT_NEAR(a, b, eps, msg) ASSERT(std::abs((a) - (b)) < (eps), msg)

int main() {
    std::cout << "Testing repeated fields..." << std::endl;

    // Test repeated scalar fields
    {
        test::repeated::RepeatedScalars msg;

        // Test initial sizes
        ASSERT_EQ(msg.bool_vals_size(), 0, "initial bool_vals size");
        ASSERT_EQ(msg.int32_vals_size(), 0, "initial int32_vals size");
        ASSERT_EQ(msg.string_vals_size(), 0, "initial string_vals size");

        // Test add and access for bool
        msg.add_bool_vals(true);
        msg.add_bool_vals(false);
        msg.add_bool_vals(true);
        ASSERT_EQ(msg.bool_vals_size(), 3, "bool_vals size after add");
        ASSERT_EQ(msg.bool_vals(0), true, "bool_vals[0]");
        ASSERT_EQ(msg.bool_vals(1), false, "bool_vals[1]");
        ASSERT_EQ(msg.bool_vals(2), true, "bool_vals[2]");

        std::cout << "  Repeated bool: OK" << std::endl;

        // Test add and access for int32
        msg.add_int32_vals(10);
        msg.add_int32_vals(20);
        msg.add_int32_vals(-30);
        ASSERT_EQ(msg.int32_vals_size(), 3, "int32_vals size");
        ASSERT_EQ(msg.int32_vals(0), 10, "int32_vals[0]");
        ASSERT_EQ(msg.int32_vals(1), 20, "int32_vals[1]");
        ASSERT_EQ(msg.int32_vals(2), -30, "int32_vals[2]");

        std::cout << "  Repeated int32: OK" << std::endl;

        // Test add and access for int64
        msg.add_int64_vals(100000000000LL);
        msg.add_int64_vals(-200000000000LL);
        ASSERT_EQ(msg.int64_vals_size(), 2, "int64_vals size");
        ASSERT_EQ(msg.int64_vals(0), 100000000000LL, "int64_vals[0]");
        ASSERT_EQ(msg.int64_vals(1), -200000000000LL, "int64_vals[1]");

        std::cout << "  Repeated int64: OK" << std::endl;

        // Test add and access for float
        msg.add_float_vals(1.5f);
        msg.add_float_vals(2.5f);
        msg.add_float_vals(3.5f);
        ASSERT_EQ(msg.float_vals_size(), 3, "float_vals size");
        ASSERT_NEAR(msg.float_vals(0), 1.5f, 0.001f, "float_vals[0]");
        ASSERT_NEAR(msg.float_vals(1), 2.5f, 0.001f, "float_vals[1]");
        ASSERT_NEAR(msg.float_vals(2), 3.5f, 0.001f, "float_vals[2]");

        std::cout << "  Repeated float: OK" << std::endl;

        // Test add and access for string
        msg.add_string_vals("hello");
        msg.add_string_vals("world");
        msg.add_string_vals("test");
        ASSERT_EQ(msg.string_vals_size(), 3, "string_vals size");
        ASSERT_EQ(msg.string_vals(0), "hello", "string_vals[0]");
        ASSERT_EQ(msg.string_vals(1), "world", "string_vals[1]");
        ASSERT_EQ(msg.string_vals(2), "test", "string_vals[2]");

        std::cout << "  Repeated string: OK" << std::endl;

        // Test clear
        msg.clear_int32_vals();
        ASSERT_EQ(msg.int32_vals_size(), 0, "int32_vals size after clear");
        // Other fields should be unaffected
        ASSERT_EQ(msg.bool_vals_size(), 3, "bool_vals size after clear_int32_vals");

        std::cout << "  Clear repeated field: OK" << std::endl;
    }

    // Test repeated message fields
    {
        test::repeated::RepeatedMessages msg;

        ASSERT_EQ(msg.items_size(), 0, "initial items size");

        // Add items
        auto* item1 = msg.add_items();
        item1->set_name("Item 1");
        item1->set_quantity(10);
        item1->set_price(9.99);

        auto* item2 = msg.add_items();
        item2->set_name("Item 2");
        item2->set_quantity(5);
        item2->set_price(19.99);

        ASSERT_EQ(msg.items_size(), 2, "items size after add");

        // Access items
        const auto& read_item1 = msg.items(0);
        ASSERT_EQ(read_item1.name(), "Item 1", "item1 name");
        ASSERT_EQ(read_item1.quantity(), 10, "item1 quantity");
        ASSERT_NEAR(read_item1.price(), 9.99, 0.001, "item1 price");

        const auto& read_item2 = msg.items(1);
        ASSERT_EQ(read_item2.name(), "Item 2", "item2 name");
        ASSERT_EQ(read_item2.quantity(), 5, "item2 quantity");
        ASSERT_NEAR(read_item2.price(), 19.99, 0.001, "item2 price");

        std::cout << "  Repeated message (add and access): OK" << std::endl;

        // Modify via mutable
        auto* mutable_item = msg.mutable_items(0);
        mutable_item->set_quantity(20);
        ASSERT_EQ(msg.items(0).quantity(), 20, "mutable modify");

        std::cout << "  Mutable repeated message: OK" << std::endl;

        // Clear
        msg.clear_items();
        ASSERT_EQ(msg.items_size(), 0, "items size after clear");

        std::cout << "  Clear repeated message: OK" << std::endl;
    }

    // Test Container with mixed fields
    {
        test::repeated::Container container;

        container.set_name("My Container");

        container.add_ids(1);
        container.add_ids(2);
        container.add_ids(3);

        container.add_tags("tag1");
        container.add_tags("tag2");

        auto* item = container.add_items();
        item->set_name("Contained Item");
        item->set_quantity(1);

        ASSERT_EQ(container.name(), "My Container", "container name");
        ASSERT_EQ(container.ids_size(), 3, "container ids size");
        ASSERT_EQ(container.tags_size(), 2, "container tags size");
        ASSERT_EQ(container.items_size(), 1, "container items size");

        std::cout << "  Mixed fields container: OK" << std::endl;
    }

    // Test direct access to RepeatedField/RepeatedPtrField
    {
        test::repeated::RepeatedScalars msg;
        msg.add_int32_vals(1);
        msg.add_int32_vals(2);
        msg.add_int32_vals(3);

        const auto& field = msg.int32_vals();
        ASSERT_EQ(field.size(), 3, "direct access size");

        std::cout << "  Direct RepeatedField access: OK" << std::endl;
    }

    std::cout << "All repeated field tests passed!" << std::endl;
    return 0;
}
