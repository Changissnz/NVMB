#include "proposition.hpp" 

void Proposition::DisplayInfo() {
    std::cout << "[x] category:\t" << category << std::endl; 
    std::cout << "[x] proposer:\t" << proposerNode << std::endl; 
    std::cout << "[x] send/recv:\t" << referenceNodes.first << "\t" << referenceNodes.second << std::endl; 
}