find_package(OpenMP)
if(OPENMP_FOUND)
set(CMAKE_C_FLAGS "${CMAKE_C_FLAGS} ${OpenMP_C_FLAGS}")
set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} ${OpenMP_CXX_FLAGS}")
set(CMAKE_EXE_LINKER_FLAGS "${CMAKE_EXE_LINKER_FLAGS} ${OpenMP_EXE_LINKER_FLAGS}")
endif()

enable_testing()
find_package(GTest REQUIRED)

# Add source directories
file(GLOB TESTFILES "*.cpp")
add_executable(testy ${TESTFILES})

target_include_directories(testy PUBLIC ${CMAKE_CURRENT_SOURCE_DIR}/../include)
target_link_libraries(testy NVMBLib ${GTEST_LIBRARIES})
