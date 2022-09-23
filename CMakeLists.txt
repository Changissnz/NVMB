cmake_minimum_required(VERSION 3.0)

project(NVMB)

include(build/cmake/GolangSetup.cmake)
add_subdirectory(golem) 
add_subdirectory(src) 
add_subdirectory(include) 
add_subdirectory(test)

# add NVMB library 
file(GLOB_RECURSE SOURCES "src/*.cpp" "include/*.h" "include/*.hpp")
add_library(NVMBLib ${SOURCES})
add_library(golemLib SHARED IMPORTED)

message("FINDING LIBRARIO: ${CMAKE_CURRENT_SOURCE_DIR}/include")
set_target_properties(golemLib PROPERTIES IMPORTED_LOCATION "${CMAKE_CURRENT_SOURCE_DIR}/include/libgolem_exec.so")

# include header folder 
target_include_directories(NVMBLib PUBLIC ${CMAKE_CURRENT_SOURCE_DIR}/test ${CMAKE_CURRENT_SOURCE_DIR}/include) 
add_executable(NVMB_EXEC "src/main.cpp")

# link the NVMBLib library and golem library to the NVMB_EXEC executable 
target_link_libraries(NVMB_EXEC PUBLIC NVMBLib)
target_link_libraries(NVMB_EXEC PUBLIC golemLib)

#target_link_libraries(NVMB_EXEC PUBLIC golem_exec) 


# use below for dynload 
target_link_libraries(NVMB_EXEC PUBLIC ${CMAKE_DL_LIBS})
