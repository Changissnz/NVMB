
set(GOPATH "${CMAKE_CURRENT_BINARY_DIR}/go")
file(MAKE_DIRECTORY ${GOPATH})

#message("CURRENT BINARY DIR : ${CMAKE_CURRENT_BINARY_DIR}")
#message("CURRENT LIST DIR : ${CMAKE_CURRENT_LIST_DIR}")

function(GO_GET TARG)
  add_custom_target(${TARG} env GOPATH=${GOPATH} go get ${ARGN})
endfunction(GO_GET)

function(ADD_GO_INSTALLABLE_PROGRAM NAME MAIN_SRC)
  get_filename_component(MAIN_SRC_ABS ${MAIN_SRC} ABSOLUTE)
  add_custom_target(${NAME})
  add_custom_command(TARGET ${NAME}
                    COMMAND env GOPATH=${GOPATH} go build 
                    -o "${CMAKE_CURRENT_BINARY_DIR}/../../include/lib${NAME}_exec.so"
                    -buildmode=c-shared ${MAIN_SRC}
                    WORKING_DIRECTORY ${CMAKE_CURRENT_LIST_DIR}
                    DEPENDS ${MAIN_SRC_ABS})
  foreach(DEP ${ARGN})
    add_dependencies(${NAME} ${DEP})
  endforeach()
  
  add_custom_target(${NAME}_all ALL DEPENDS ${NAME})
  install(PROGRAMS ${CMAKE_CURRENT_BINARY_DIR}/${NAME} DESTINATION bin)
endfunction(ADD_GO_INSTALLABLE_PROGRAM)