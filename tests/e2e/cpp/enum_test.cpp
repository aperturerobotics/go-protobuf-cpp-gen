// Enum types test - verifies enum handling works correctly
#include <iostream>
#include <string>

#include "enums.pb.h"

#define ASSERT(cond, msg) \
    if (!(cond)) { \
        std::cerr << "FAILED: " << msg << " at " << __FILE__ << ":" << __LINE__ << std::endl; \
        return 1; \
    }

#define ASSERT_EQ(a, b, msg) ASSERT((a) == (b), msg)
#define ASSERT_STREQ(a, b, msg) ASSERT(std::string(a) == std::string(b), msg)

int main() {
    std::cout << "Testing enum types..." << std::endl;

    // Test Status enum values
    {
        ASSERT_EQ(test::enums::STATUS_UNKNOWN, 0, "STATUS_UNKNOWN value");
        ASSERT_EQ(test::enums::STATUS_ACTIVE, 1, "STATUS_ACTIVE value");
        ASSERT_EQ(test::enums::STATUS_INACTIVE, 2, "STATUS_INACTIVE value");
        ASSERT_EQ(test::enums::STATUS_PENDING, 3, "STATUS_PENDING value");

        std::cout << "  Status enum values: OK" << std::endl;
    }

    // Test Priority enum values (with gaps)
    {
        ASSERT_EQ(test::enums::PRIORITY_UNSPECIFIED, 0, "PRIORITY_UNSPECIFIED value");
        ASSERT_EQ(test::enums::PRIORITY_LOW, 10, "PRIORITY_LOW value");
        ASSERT_EQ(test::enums::PRIORITY_MEDIUM, 20, "PRIORITY_MEDIUM value");
        ASSERT_EQ(test::enums::PRIORITY_HIGH, 30, "PRIORITY_HIGH value");
        ASSERT_EQ(test::enums::PRIORITY_CRITICAL, 100, "PRIORITY_CRITICAL value");

        std::cout << "  Priority enum values: OK" << std::endl;
    }

    // Test Status_IsValid
    {
        ASSERT(test::enums::Status_IsValid(0), "Status_IsValid(0)");
        ASSERT(test::enums::Status_IsValid(1), "Status_IsValid(1)");
        ASSERT(test::enums::Status_IsValid(2), "Status_IsValid(2)");
        ASSERT(test::enums::Status_IsValid(3), "Status_IsValid(3)");
        ASSERT(!test::enums::Status_IsValid(4), "Status_IsValid(4) should be false");
        ASSERT(!test::enums::Status_IsValid(-1), "Status_IsValid(-1) should be false");
        ASSERT(!test::enums::Status_IsValid(100), "Status_IsValid(100) should be false");

        std::cout << "  Status_IsValid: OK" << std::endl;
    }

    // Test Priority_IsValid
    {
        ASSERT(test::enums::Priority_IsValid(0), "Priority_IsValid(0)");
        ASSERT(test::enums::Priority_IsValid(10), "Priority_IsValid(10)");
        ASSERT(test::enums::Priority_IsValid(20), "Priority_IsValid(20)");
        ASSERT(test::enums::Priority_IsValid(30), "Priority_IsValid(30)");
        ASSERT(test::enums::Priority_IsValid(100), "Priority_IsValid(100)");
        ASSERT(!test::enums::Priority_IsValid(5), "Priority_IsValid(5) should be false");
        ASSERT(!test::enums::Priority_IsValid(15), "Priority_IsValid(15) should be false");

        std::cout << "  Priority_IsValid: OK" << std::endl;
    }

    // Test Status_Name
    {
        ASSERT_STREQ(test::enums::Status_Name(test::enums::STATUS_UNKNOWN), "STATUS_UNKNOWN", "Status_Name UNKNOWN");
        ASSERT_STREQ(test::enums::Status_Name(test::enums::STATUS_ACTIVE), "STATUS_ACTIVE", "Status_Name ACTIVE");
        ASSERT_STREQ(test::enums::Status_Name(test::enums::STATUS_INACTIVE), "STATUS_INACTIVE", "Status_Name INACTIVE");
        ASSERT_STREQ(test::enums::Status_Name(test::enums::STATUS_PENDING), "STATUS_PENDING", "Status_Name PENDING");

        std::cout << "  Status_Name: OK" << std::endl;
    }

    // Test enum MIN/MAX
    {
        ASSERT_EQ(test::enums::Status_MIN, test::enums::STATUS_UNKNOWN, "Status_MIN");
        ASSERT_EQ(test::enums::Status_MAX, test::enums::STATUS_PENDING, "Status_MAX");

        std::cout << "  Enum MIN/MAX: OK" << std::endl;
    }

    // Test EnumMessage
    {
        test::enums::EnumMessage msg;

        // Test default values
        ASSERT_EQ(msg.status(), test::enums::STATUS_UNKNOWN, "default status");
        ASSERT_EQ(msg.priority(), test::enums::PRIORITY_UNSPECIFIED, "default priority");

        // Test setters
        msg.set_status(test::enums::STATUS_ACTIVE);
        msg.set_priority(test::enums::PRIORITY_HIGH);

        ASSERT_EQ(msg.status(), test::enums::STATUS_ACTIVE, "status after set");
        ASSERT_EQ(msg.priority(), test::enums::PRIORITY_HIGH, "priority after set");

        // Test clear
        msg.clear_status();
        ASSERT_EQ(msg.status(), test::enums::STATUS_UNKNOWN, "status after clear");

        std::cout << "  EnumMessage accessors: OK" << std::endl;
    }

    // Test repeated enum field
    {
        test::enums::EnumMessage msg;

        ASSERT_EQ(msg.statuses_size(), 0, "initial statuses size");

        msg.add_statuses(test::enums::STATUS_ACTIVE);
        msg.add_statuses(test::enums::STATUS_PENDING);
        msg.add_statuses(test::enums::STATUS_INACTIVE);

        ASSERT_EQ(msg.statuses_size(), 3, "statuses size after add");
        ASSERT_EQ(msg.statuses(0), test::enums::STATUS_ACTIVE, "statuses[0]");
        ASSERT_EQ(msg.statuses(1), test::enums::STATUS_PENDING, "statuses[1]");
        ASSERT_EQ(msg.statuses(2), test::enums::STATUS_INACTIVE, "statuses[2]");

        msg.clear_statuses();
        ASSERT_EQ(msg.statuses_size(), 0, "statuses size after clear");

        std::cout << "  Repeated enum field: OK" << std::endl;
    }

    // Test nested enum
    {
        test::enums::MessageWithNestedEnum msg;

        ASSERT_EQ(msg.status(), test::enums::NESTED_UNKNOWN, "default nested status");

        msg.set_status(test::enums::NESTED_OK);
        ASSERT_EQ(msg.status(), test::enums::NESTED_OK, "nested status after set");

        msg.set_message("test message");
        ASSERT_EQ(msg.message(), "test message", "nested message string");

        std::cout << "  Nested enum: OK" << std::endl;
    }

    std::cout << "All enum tests passed!" << std::endl;
    return 0;
}
