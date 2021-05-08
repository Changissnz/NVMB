#ifndef NVMB_NODE_HPP
#define NVMB_NODE_HPP

#include <string>
#include <vector>
#include <functional> 
#include <armadillo>
using namespace std; 
using namespace arma;

vector<string> tvcf(string); 

class Vuln
{
public: 

    string identifier; 
    float quality; 
    float span; 
    function<vector<string>(string)> vcf; // patch function; return, (args)
    Vuln(string idf, float q, float s, vector<string> (*f)(string));

    // TODO: delete after finish, test method for Armadillo
    void CalculatePlans();
    
}; 

#endif