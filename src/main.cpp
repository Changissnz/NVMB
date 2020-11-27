#include <iostream>
#include <dlfcn.h>
#include <stdio.h> 
#include <stdlib.h>
#include "libgolem_exec.h"
#include "network_generation.hpp" 

/// uncomment below for dynamic-loading 
//#include <unistd.h>
using namespace std; 

/*
types of modifications: 

- need to add the strategy moves node chooses each round. 
*/ 
int main(int argc, char **argv) {
	cout << "HELLLO" << endl; 
	cout << "CALLING GO FUNC " << endl;

    GoInt a = 12;
    GoInt b = 99;
    printf("awesome.Add(12,99) = %d %d\n", Add(a, b), IncorpGolem());

    cout << "incorporating golem = " << IncorpGolem() << endl; 
	return -1; 
} 
