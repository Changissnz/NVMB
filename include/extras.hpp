#ifndef EXTRAS_HPP
#define EXTRAS_HPP

#include <string>
#include <map>
#include <set>
#include <vector>
#include <utility>
#include <algorithm> 
#include <iterator>
#include <iostream>
#include <assert.h>

std::string NodePairToString(int v1, int v2);
int StringToInt(std::string s);
std::pair<int,int> StringToNodePair(std::string E);

template <typename T>
std::set<T> UnionSet(std::set<T> v1, std::set<T> v2) // #2
{
    std::set<T> v3;
    typename std::set<T>::iterator it;
    for (it = v1.begin(); it != v1.end(); it++) {
        v3.insert(*it);
    }

    for (it = v2.begin(); it != v2.end(); it++) {
        v3.insert(*it);
    }

    return v3;
};

template <typename T>
std::set<T> RemoveSet(std::set<T> v1, std::set<T> v2) {
    for (auto it = v2.begin(); it != v2.end(); it++) {
        v1.erase(*it);
    };
    return v1;
};

template <typename T> 
std::vector<T> SetToShuffledVector(std::set<T> s1) { 
    std::vector<T> vec(s1.begin(), s1.end()); 
    std::random_shuffle ( vec.begin(), vec.end()); 
    return vec; 
}

/// CAUTION: method is used only for basic types
template<typename C>
void DisplayIterable(const C& container)
{
    for(auto v: container)
        std::cout << v << " ";
    std::cout << std::endl;
};

/// CAUTION: use only with simple iterables. 
template<typename C>  
std::string NumIterToStringType(const C& container) {
    std::string output = ""; 
    for(auto v: container)
      output += std::to_string(v) + " "; 

    return output; 
}

template<typename C1,typename C2> 
void DisplayMap(std::map<C1,C2> m) {
  for (auto x: m) {
    std::cout << x.first << std::endl; 
    std::cout << x.second << std::endl; 
    std::cout << "\n" << std::endl; 
  }
}

float ZeroDiv(float num, float denum, float defaultVal);
float RandomFloat(); 

#endif