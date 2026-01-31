// Map fields test - verifies map field handling
#include <iostream>
#include <string>
#include <cmath>

#include "maps.pb.h"

#define ASSERT(cond, msg) \
    if (!(cond)) { \
        std::cerr << "FAILED: " << msg << " at " << __FILE__ << ":" << __LINE__ << std::endl; \
        return 1; \
    }

#define ASSERT_EQ(a, b, msg) ASSERT((a) == (b), msg)
#define ASSERT_NEAR(a, b, eps, msg) ASSERT(std::abs((a) - (b)) < (eps), msg)

int main() {
    std::cout << "Testing map fields..." << std::endl;

    // Test string-to-scalar maps
    {
        test::maps::MapScalars msg;

        // Test initial size
        ASSERT_EQ(msg.string_to_string_size(), 0, "initial string_to_string size");
        ASSERT_EQ(msg.string_to_int32_size(), 0, "initial string_to_int32 size");

        // Add to string->string map
        (*msg.mutable_string_to_string())["key1"] = "value1";
        (*msg.mutable_string_to_string())["key2"] = "value2";

        ASSERT_EQ(msg.string_to_string_size(), 2, "string_to_string size");

        // Access via const accessor
        const auto& str_map = msg.string_to_string();
        auto it1 = str_map.find("key1");
        ASSERT(it1 != str_map.end(), "key1 found");
        ASSERT_EQ(it1->second, "value1", "key1 value");

        auto it2 = str_map.find("key2");
        ASSERT(it2 != str_map.end(), "key2 found");
        ASSERT_EQ(it2->second, "value2", "key2 value");

        std::cout << "  String->String map: OK" << std::endl;

        // Add to string->int32 map
        (*msg.mutable_string_to_int32())["score1"] = 100;
        (*msg.mutable_string_to_int32())["score2"] = 200;

        ASSERT_EQ(msg.string_to_int32_size(), 2, "string_to_int32 size");
        ASSERT_EQ(msg.string_to_int32().at("score1"), 100, "score1 value");
        ASSERT_EQ(msg.string_to_int32().at("score2"), 200, "score2 value");

        std::cout << "  String->Int32 map: OK" << std::endl;

        // Add to string->bool map
        (*msg.mutable_string_to_bool())["flag1"] = true;
        (*msg.mutable_string_to_bool())["flag2"] = false;

        ASSERT_EQ(msg.string_to_bool().at("flag1"), true, "flag1 value");
        ASSERT_EQ(msg.string_to_bool().at("flag2"), false, "flag2 value");

        std::cout << "  String->Bool map: OK" << std::endl;

        // Add to string->double map
        (*msg.mutable_string_to_double())["pi"] = 3.14159;
        (*msg.mutable_string_to_double())["e"] = 2.71828;

        ASSERT_NEAR(msg.string_to_double().at("pi"), 3.14159, 0.00001, "pi value");
        ASSERT_NEAR(msg.string_to_double().at("e"), 2.71828, 0.00001, "e value");

        std::cout << "  String->Double map: OK" << std::endl;
    }

    // Test int-keyed maps
    {
        test::maps::MapScalars msg;

        (*msg.mutable_int32_to_string())[1] = "one";
        (*msg.mutable_int32_to_string())[2] = "two";
        (*msg.mutable_int32_to_string())[-1] = "negative one";

        ASSERT_EQ(msg.int32_to_string_size(), 3, "int32_to_string size");
        ASSERT_EQ(msg.int32_to_string().at(1), "one", "1 -> one");
        ASSERT_EQ(msg.int32_to_string().at(2), "two", "2 -> two");
        ASSERT_EQ(msg.int32_to_string().at(-1), "negative one", "-1 -> negative one");

        std::cout << "  Int32->String map: OK" << std::endl;

        (*msg.mutable_int64_to_string())[1000000000000LL] = "trillion";
        (*msg.mutable_int64_to_string())[-1000000000000LL] = "negative trillion";

        ASSERT_EQ(msg.int64_to_string().at(1000000000000LL), "trillion", "int64 key 1");
        ASSERT_EQ(msg.int64_to_string().at(-1000000000000LL), "negative trillion", "int64 key 2");

        std::cout << "  Int64->String map: OK" << std::endl;

        (*msg.mutable_bool_to_string())[true] = "yes";
        (*msg.mutable_bool_to_string())[false] = "no";

        ASSERT_EQ(msg.bool_to_string().at(true), "yes", "true -> yes");
        ASSERT_EQ(msg.bool_to_string().at(false), "no", "false -> no");

        std::cout << "  Bool->String map: OK" << std::endl;
    }

    // Test maps with message values
    {
        test::maps::MapMessages msg;

        // Add message values
        test::maps::MapValue val1;
        val1.set_name("value1");
        val1.set_score(100);
        (*msg.mutable_string_to_message())["key1"] = val1;

        test::maps::MapValue val2;
        val2.set_name("value2");
        val2.set_score(200);
        (*msg.mutable_string_to_message())["key2"] = val2;

        ASSERT_EQ(msg.string_to_message_size(), 2, "string_to_message size");

        const auto& retrieved1 = msg.string_to_message().at("key1");
        ASSERT_EQ(retrieved1.name(), "value1", "retrieved1 name");
        ASSERT_EQ(retrieved1.score(), 100, "retrieved1 score");

        const auto& retrieved2 = msg.string_to_message().at("key2");
        ASSERT_EQ(retrieved2.name(), "value2", "retrieved2 name");
        ASSERT_EQ(retrieved2.score(), 200, "retrieved2 score");

        std::cout << "  String->Message map: OK" << std::endl;

        // Test int32->message map
        test::maps::MapValue val3;
        val3.set_name("value3");
        val3.set_score(300);
        (*msg.mutable_int32_to_message())[42] = val3;

        ASSERT_EQ(msg.int32_to_message().at(42).name(), "value3", "int32->message name");
        ASSERT_EQ(msg.int32_to_message().at(42).score(), 300, "int32->message score");

        std::cout << "  Int32->Message map: OK" << std::endl;
    }

    // Test mixed fields with maps
    {
        test::maps::MixedWithMaps msg;

        msg.set_name("mixed message");

        (*msg.mutable_scores())["alice"] = 95;
        (*msg.mutable_scores())["bob"] = 87;

        msg.add_tags("important");
        msg.add_tags("review");

        test::maps::MapValue data1;
        data1.set_name("data1");
        data1.set_score(50);
        (*msg.mutable_data())["entry1"] = data1;

        ASSERT_EQ(msg.name(), "mixed message", "mixed name");
        ASSERT_EQ(msg.scores_size(), 2, "mixed scores size");
        ASSERT_EQ(msg.scores().at("alice"), 95, "mixed alice score");
        ASSERT_EQ(msg.tags_size(), 2, "mixed tags size");
        ASSERT_EQ(msg.data_size(), 1, "mixed data size");

        std::cout << "  Mixed fields with maps: OK" << std::endl;
    }

    // Test map clear
    {
        test::maps::MapScalars msg;

        (*msg.mutable_string_to_string())["key"] = "value";
        ASSERT_EQ(msg.string_to_string_size(), 1, "before clear");

        msg.clear_string_to_string();
        ASSERT_EQ(msg.string_to_string_size(), 0, "after clear");

        std::cout << "  Map clear: OK" << std::endl;
    }

    // Test map iteration
    {
        test::maps::MapScalars msg;

        (*msg.mutable_string_to_int32())["a"] = 1;
        (*msg.mutable_string_to_int32())["b"] = 2;
        (*msg.mutable_string_to_int32())["c"] = 3;

        int sum = 0;
        for (const auto& pair : msg.string_to_int32()) {
            sum += pair.second;
        }
        ASSERT_EQ(sum, 6, "iteration sum");

        std::cout << "  Map iteration: OK" << std::endl;
    }

    // Test map copy
    {
        test::maps::MapScalars msg1;
        (*msg1.mutable_string_to_string())["key"] = "value";

        test::maps::MapScalars msg2(msg1);
        ASSERT_EQ(msg2.string_to_string().at("key"), "value", "copy map value");

        // Modify copy
        (*msg2.mutable_string_to_string())["key"] = "modified";
        ASSERT_EQ(msg1.string_to_string().at("key"), "value", "original unchanged");
        ASSERT_EQ(msg2.string_to_string().at("key"), "modified", "copy modified");

        std::cout << "  Map copy: OK" << std::endl;
    }

    std::cout << "All map tests passed!" << std::endl;
    return 0;
}
