cmake_minimum_required(VERSION 3.14)
project(protobuf_cpp_gen_tests CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)

# Find protobuf - try multiple methods
find_package(Protobuf QUIET)
if(NOT Protobuf_FOUND)
    # Try with hints for homebrew on macOS
    if(APPLE)
        execute_process(
            COMMAND brew --prefix protobuf
            OUTPUT_VARIABLE PROTOBUF_PREFIX
            OUTPUT_STRIP_TRAILING_WHITESPACE
            ERROR_QUIET
        )
        if(PROTOBUF_PREFIX)
            list(APPEND CMAKE_PREFIX_PATH "${PROTOBUF_PREFIX}")
            find_package(Protobuf QUIET)
        endif()
    endif()
endif()

if(NOT Protobuf_FOUND)
    message(FATAL_ERROR "Protobuf not found. Please install protobuf library.")
endif()

message(STATUS "Protobuf version: ${Protobuf_VERSION}")
message(STATUS "Protobuf include: ${Protobuf_INCLUDE_DIRS}")
message(STATUS "Protobuf libraries: ${Protobuf_LIBRARIES}")

# Include directories
include_directories(${Protobuf_INCLUDE_DIRS})
include_directories(${CMAKE_CURRENT_SOURCE_DIR}/generated)

# Generated source files
set(GENERATED_SOURCES
    generated/scalars.pb.cc
    generated/enums.pb.cc
    generated/repeated.pb.cc
    generated/nested.pb.cc
    generated/maps.pb.cc
    generated/common.pb.cc
    generated/imports.pb.cc
)

set(GENERATED_HEADERS
    generated/scalars.pb.h
    generated/enums.pb.h
    generated/repeated.pb.h
    generated/nested.pb.h
    generated/maps.pb.h
    generated/common.pb.h
    generated/imports.pb.h
)

# Create a library from generated sources
add_library(test_protos STATIC ${GENERATED_SOURCES})
target_link_libraries(test_protos ${Protobuf_LIBRARIES})
target_include_directories(test_protos PUBLIC ${CMAKE_CURRENT_SOURCE_DIR}/generated)

# Suppress warnings in generated code
target_compile_options(test_protos PRIVATE -w)

# Test executable for compilation verification
add_executable(compile_test cpp/compile_test.cpp)
target_link_libraries(compile_test test_protos ${Protobuf_LIBRARIES})

# Test executable for scalar types
add_executable(scalar_test cpp/scalar_test.cpp)
target_link_libraries(scalar_test test_protos ${Protobuf_LIBRARIES})

# Test executable for enums
add_executable(enum_test cpp/enum_test.cpp)
target_link_libraries(enum_test test_protos ${Protobuf_LIBRARIES})

# Test executable for repeated fields
add_executable(repeated_test cpp/repeated_test.cpp)
target_link_libraries(repeated_test test_protos ${Protobuf_LIBRARIES})

# Test executable for nested messages
add_executable(nested_test cpp/nested_test.cpp)
target_link_libraries(nested_test test_protos ${Protobuf_LIBRARIES})

# Test executable for maps
add_executable(map_test cpp/map_test.cpp)
target_link_libraries(map_test test_protos ${Protobuf_LIBRARIES})

# Test executable for serialization round-trip
add_executable(serialize_test cpp/serialize_test.cpp)
target_link_libraries(serialize_test test_protos ${Protobuf_LIBRARIES})

# Enable testing
enable_testing()

add_test(NAME compile_test COMMAND compile_test)
add_test(NAME scalar_test COMMAND scalar_test)
add_test(NAME enum_test COMMAND enum_test)
add_test(NAME repeated_test COMMAND repeated_test)
add_test(NAME nested_test COMMAND nested_test)
add_test(NAME map_test COMMAND map_test)
add_test(NAME serialize_test COMMAND serialize_test)
