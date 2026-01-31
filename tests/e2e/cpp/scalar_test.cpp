// Scalar types test - verifies getters/setters work correctly
#include <iostream>
#include <string>
#include <cmath>
#include <cstdlib>

#include "scalars.pb.h"

#define ASSERT(cond, msg) \
    if (!(cond)) { \
        std::cerr << "FAILED: " << msg << " at " << __FILE__ << ":" << __LINE__ << std::endl; \
        return 1; \
    }

#define ASSERT_EQ(a, b, msg) ASSERT((a) == (b), msg)
#define ASSERT_NEAR(a, b, eps, msg) ASSERT(std::abs((a) - (b)) < (eps), msg)

int main() {
    std::cout << "Testing scalar field accessors..." << std::endl;

    // Test AllScalars message
    {
        test::scalars::AllScalars msg;

        // Test default values (proto3 defaults to zero/empty)
        ASSERT_EQ(msg.bool_val(), false, "default bool should be false");
        ASSERT_EQ(msg.int32_val(), 0, "default int32 should be 0");
        ASSERT_EQ(msg.int64_val(), 0, "default int64 should be 0");
        ASSERT_EQ(msg.uint32_val(), 0u, "default uint32 should be 0");
        ASSERT_EQ(msg.uint64_val(), 0u, "default uint64 should be 0");
        ASSERT_EQ(msg.sint32_val(), 0, "default sint32 should be 0");
        ASSERT_EQ(msg.sint64_val(), 0, "default sint64 should be 0");
        ASSERT_EQ(msg.fixed32_val(), 0u, "default fixed32 should be 0");
        ASSERT_EQ(msg.fixed64_val(), 0u, "default fixed64 should be 0");
        ASSERT_EQ(msg.sfixed32_val(), 0, "default sfixed32 should be 0");
        ASSERT_EQ(msg.sfixed64_val(), 0, "default sfixed64 should be 0");
        ASSERT_NEAR(msg.float_val(), 0.0f, 0.0001f, "default float should be 0.0");
        ASSERT_NEAR(msg.double_val(), 0.0, 0.0001, "default double should be 0.0");
        ASSERT_EQ(msg.string_val(), "", "default string should be empty");
        ASSERT_EQ(msg.bytes_val(), "", "default bytes should be empty");

        std::cout << "  Default values: OK" << std::endl;

        // Test setters
        msg.set_bool_val(true);
        msg.set_int32_val(-12345);
        msg.set_int64_val(-9876543210LL);
        msg.set_uint32_val(12345u);
        msg.set_uint64_val(9876543210ULL);
        msg.set_sint32_val(-54321);
        msg.set_sint64_val(-1234567890LL);
        msg.set_fixed32_val(0xDEADBEEFu);
        msg.set_fixed64_val(0xDEADBEEFCAFEBABEULL);
        msg.set_sfixed32_val(-0x12345678);
        msg.set_sfixed64_val(-0x123456789ABCDEFLL);
        msg.set_float_val(3.14159f);
        msg.set_double_val(2.718281828);
        msg.set_string_val("hello world");
        msg.set_bytes_val("binary\x00data");

        // Verify getters return set values
        ASSERT_EQ(msg.bool_val(), true, "bool getter");
        ASSERT_EQ(msg.int32_val(), -12345, "int32 getter");
        ASSERT_EQ(msg.int64_val(), -9876543210LL, "int64 getter");
        ASSERT_EQ(msg.uint32_val(), 12345u, "uint32 getter");
        ASSERT_EQ(msg.uint64_val(), 9876543210ULL, "uint64 getter");
        ASSERT_EQ(msg.sint32_val(), -54321, "sint32 getter");
        ASSERT_EQ(msg.sint64_val(), -1234567890LL, "sint64 getter");
        ASSERT_EQ(msg.fixed32_val(), 0xDEADBEEFu, "fixed32 getter");
        ASSERT_EQ(msg.fixed64_val(), 0xDEADBEEFCAFEBABEULL, "fixed64 getter");
        ASSERT_EQ(msg.sfixed32_val(), -0x12345678, "sfixed32 getter");
        ASSERT_EQ(msg.sfixed64_val(), -0x123456789ABCDEFLL, "sfixed64 getter");
        ASSERT_NEAR(msg.float_val(), 3.14159f, 0.0001f, "float getter");
        ASSERT_NEAR(msg.double_val(), 2.718281828, 0.0001, "double getter");
        ASSERT_EQ(msg.string_val(), "hello world", "string getter");

        std::cout << "  Setters and getters: OK" << std::endl;

        // Test clear methods
        msg.clear_bool_val();
        msg.clear_int32_val();
        msg.clear_string_val();

        ASSERT_EQ(msg.bool_val(), false, "cleared bool");
        ASSERT_EQ(msg.int32_val(), 0, "cleared int32");
        ASSERT_EQ(msg.string_val(), "", "cleared string");

        std::cout << "  Clear methods: OK" << std::endl;

        // Test mutable string access
        std::string* mutable_str = msg.mutable_string_val();
        *mutable_str = "modified";
        ASSERT_EQ(msg.string_val(), "modified", "mutable string");

        std::cout << "  Mutable accessors: OK" << std::endl;

        // Test move semantics for string
        std::string temp = "moved value";
        msg.set_string_val(std::move(temp));
        ASSERT_EQ(msg.string_val(), "moved value", "move set string");

        std::cout << "  Move semantics: OK" << std::endl;
    }

    // Test Clear() method on whole message
    {
        test::scalars::AllScalars msg;
        msg.set_bool_val(true);
        msg.set_int32_val(42);
        msg.set_string_val("test");

        msg.Clear();

        ASSERT_EQ(msg.bool_val(), false, "Clear() bool");
        ASSERT_EQ(msg.int32_val(), 0, "Clear() int32");
        ASSERT_EQ(msg.string_val(), "", "Clear() string");

        std::cout << "  Clear() method: OK" << std::endl;
    }

    // Test copy constructor and assignment
    {
        test::scalars::AllScalars msg1;
        msg1.set_int32_val(100);
        msg1.set_string_val("original");

        test::scalars::AllScalars msg2(msg1);
        ASSERT_EQ(msg2.int32_val(), 100, "copy constructor int32");
        ASSERT_EQ(msg2.string_val(), "original", "copy constructor string");

        test::scalars::AllScalars msg3;
        msg3 = msg1;
        ASSERT_EQ(msg3.int32_val(), 100, "copy assignment int32");
        ASSERT_EQ(msg3.string_val(), "original", "copy assignment string");

        std::cout << "  Copy semantics: OK" << std::endl;
    }

    // Test move constructor and assignment
    {
        test::scalars::AllScalars msg1;
        msg1.set_int32_val(200);
        msg1.set_string_val("to be moved");

        test::scalars::AllScalars msg2(std::move(msg1));
        ASSERT_EQ(msg2.int32_val(), 200, "move constructor int32");
        ASSERT_EQ(msg2.string_val(), "to be moved", "move constructor string");

        test::scalars::AllScalars msg3;
        msg3.set_int32_val(300);
        msg3.set_string_val("also to be moved");

        test::scalars::AllScalars msg4;
        msg4 = std::move(msg3);
        ASSERT_EQ(msg4.int32_val(), 300, "move assignment int32");
        ASSERT_EQ(msg4.string_val(), "also to be moved", "move assignment string");

        std::cout << "  Move semantics: OK" << std::endl;
    }

    std::cout << "All scalar tests passed!" << std::endl;
    return 0;
}
