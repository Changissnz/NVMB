#include "nvmb_node.hpp"

/*
constructors for Node class 
*/ 

NVMBNode::NVMBNode(int id, string fn) {
    identifier = id;
    reader = new NodeFileReader(fn);
    NodeStrategos* strategy = NodeStrategos::MakeOne(map<string, string>());
    npu = new NodeProcessingUnit(identifier, set<int>(), strategy);
    bank = new BankUnit();
    mode = "active";
    knownNodes = set<int>(); 
    testModeSpec = ""; 
    manInput = false; 
    initialized = false;
    timestamp = 0; 
    discoveryCost = 0;
    defaultResponseCF = make_pair(false,true);
    neighborIdentifiers = set<int>(); 
}

NVMBNode::NVMBNode(int id, string fn, string m) {
    identifier = id;
    reader = new NodeFileReader(fn);
    NodeStrategos* strategy = NodeStrategos::MakeOne(map<string, string>());
    npu = new NodeProcessingUnit(identifier, set<int>(), strategy);
    bank = new BankUnit();
    assert (m == "active" || m == "inactive");
    mode = m;
    knownNodes = set<int>(); 
    testModeSpec = ""; 
    manInput = false; 
    initialized = false;
    timestamp = 0; 
    discoveryCost = 0;
    defaultResponseCF = make_pair(false,true);
    neighborIdentifiers = set<int>(); 
}

NVMBNode::NVMBNode(int id, string fn, map<string,string> m) {
    identifier = id;
    reader = new NodeFileReader(fn);
    NodeStrategos* strategy = NodeStrategos::MakeOne(m);
    npu = new NodeProcessingUnit(identifier, set<int>(), strategy);
    bank = new BankUnit();
    mode = "active";
    knownNodes = set<int>(); 
    testModeSpec = ""; 
    manInput = false; 
    initialized = false;
    timestamp = 0; 
    discoveryCost = 0;
    defaultResponseCF = make_pair(false,true);
    neighborIdentifiers = set<int>();
}

NVMBNode::NVMBNode(int id, string dd, string fn, map<string,string> m) {
    identifier = id; 
    reader = new NodeFileReader(dd, fn);
    NodeStrategos* strategy = NodeStrategos::MakeOne(m);
    npu = new NodeProcessingUnit(identifier, set<int>(), strategy);
    bank = new BankUnit();
    mode = "active";
    testModeSpec = ""; 
    manInput = false; 
    initialized = false;
    discoveryCost = 0;
    timestamp = 0; 
    defaultResponseCF = make_pair(false,true);
    neighborIdentifiers = set<int>(); 
}