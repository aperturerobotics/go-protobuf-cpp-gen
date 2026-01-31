// Compilation test - verifies all generated headers compile correctly
#include <iostream>
#include <string>

#include "scalars.pb.h"
#include "enums.pb.h"
#include "repeated.pb.h"
#include "nested.pb.h"
#include "maps.pb.h"
#include "common.pb.h"
#include "imports.pb.h"

int main() {
    std::cout << "Testing compilation of generated protobuf code..." << std::endl;

    // Instantiate all message types to verify they compile
    {
        test::scalars::AllScalars msg;
        test::scalars::DefaultValues defaults;
        std::cout << "  scalars.proto: OK" << std::endl;
    }

    {
        test::enums::EnumMessage msg;
        test::enums::MessageWithNestedEnum nested;
        std::cout << "  enums.proto: OK" << std::endl;
    }

    {
        test::repeated::RepeatedScalars scalars;
        test::repeated::Item item;
        test::repeated::RepeatedMessages msgs;
        test::repeated::Container container;
        std::cout << "  repeated.proto: OK" << std::endl;
    }

    {
        test::nested::Outer outer;
        test::nested::Middle middle;  // All types at namespace level
        test::nested::Inner inner;
        test::nested::Document doc;
        test::nested::Level1 level1;
        std::cout << "  nested.proto: OK" << std::endl;
    }

    {
        test::maps::MapScalars scalars;
        test::maps::MapValue value;
        test::maps::MapMessages msgs;
        test::maps::MixedWithMaps mixed;
        std::cout << "  maps.proto: OK" << std::endl;
    }

    {
        test::common::Timestamp ts;
        test::common::Duration dur;
        test::common::Metadata meta;
        std::cout << "  common.proto: OK" << std::endl;
    }

    {
        test::imports::Resource res;
        test::imports::Task task;
        test::imports::Event event;
        std::cout << "  imports.proto: OK" << std::endl;
    }

    std::cout << "All compilation tests passed!" << std::endl;
    return 0;
}
