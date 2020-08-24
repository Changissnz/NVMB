// network contains all node instances and represents the network of virtual
// machines. Methods for network have been

#include "network.hpp"
using namespace std;

int MAX_FLARES = 5000; 
int MAX_FLARE_MULTIPLIER = 4;

/*
sets the node nature variables to random variables 
*/ 
map<string,string> MakeRandomNodeNatureVariables() {
    map<string,string> natureVar = map<string,string>();
    natureVar["competition"] = to_string(RandomFloat()); 
    natureVar["greed"] = to_string(RandomFloat()); 
    natureVar["negotiation"] = to_string(RandomFloat()); 
    natureVar["growth"] = to_string(RandomFloat()); 
    return natureVar; 
}

bool Network::NodeExists(int nodeIdentifier) {
  return (contents.find(nodeIdentifier) == contents.end()) ? false : true;
}

/*
// description
adds node by identifier, will auto-create Node instance
by default file `node_i`
*/
bool Network::AddNode(int i) {
  if (contents.find(i) == contents.end() ) {
    char buffer [50];
    sprintf (buffer, "node_%d", i);
    NVMBNode n(i, buffer);
    contents[i] = &n;
    return true;
  }
  return false;
}

/*
// description
adds node by Node instance
*/
bool Network::AddNode(NVMBNode* n) {
  if (contents.find(n->GetId()) == contents.end() ) {
    contents[n->GetId()] = n;
    return true;
  }
  return false;
}

/*
// description
for a given edge (x1,x2), breaks the edge between them.
If directedBreak is -1, breaks both edges,
*/
void Network::BreakEdge(int v1, int v2, int directedBreak) {
    if (contents[v1] != nullptr) {
        contents[v1]->DeleteNeighbor(v2);
    }

    if (!directedBreak && contents[v2] != nullptr) {
        contents[v2]->DeleteNeighbor(v1);
    }
}

/*
// description
Adds an edge to the network. If `directedAdd` is set to true, then
adds the edge v1-v2, otherwise adds the two edges v1-v2 and v2-v1.

CAUTION : no error-handling

// arguments
v1 - first node identifier
v2 - second node identifier
directedAdd - specifies if edge is directed or not
*/
void Network::AddEdge(int v1, int v2, int directedAdd) {
    contents[v1]->AddNeighbor(v2);
    if (!directedAdd) {
        contents[v2]->AddNeighbor(v1);
    }
}

/*
// description
Checks if edge exists

CAUTION: no argument error-handling
*/
bool Network::EdgeExists(int v1, int v2, int directed) {

    if (contents.find(v1) == contents.end() || contents.find(v2) == contents.end()) {
        return false; 
    }

    if (directed) {
        return (contents[v1]->NeighborExists(v2)) ? true : false;
    }

    return (contents[v1]->NeighborExists(v2) && contents[v2]->NeighborExists(v1)) ? true : false;
}

/*
// description
gets # of nodes of network
*/
int Network::Size() {
  return contents.size();
}

/*
// description
for each node,
- update its timestamp
- log its bank's timestamp data
*/
void Network::UpdateNodeData() {

    timestamp++; 
    NVMBNode* n;
    BankUnit* b;

    //#pragma omp parallel num_threads(DEFAULT_NUM_THREADS)
    //{
    //    #pragma omp for
        for (map<int, NVMBNode*>::iterator it = contents.begin(); it != contents.end(); it++) {
            n = it->second;
            if (n == nullptr) {continue;} 
            n->LogTimestampAttribute(); 
            b = n->GetBank();
            b->LogTimestampData();
            n->SetTimestamp(timestamp); 
        };

    return;
}

/// TODO: delete this.
void Network::UpdateTime() {
    timestamp++; 
    for (auto c:contents) {
        if (c.second != nullptr) {
            auto b = (c.second)->GetBank();
            b->LogTimestampData();
            (c.second)->SetTimestamp(timestamp);  
        }
    }
}

/*
// description
step 1: updates each node's owned
- each owned flare will update itself with its cache

step 2: assign each held flare to its correct node location
*/
void Network::UpdateNodes(int verbose) {

    // allocate held DF to nodes (current location) using DF info.

    UpdateNodeDFOwned();
    AssignHeldDFAll();
    UpdateNodeDFNeighbors();

    // update commerce flare
    if (verbose == 1 || verbose == 2) {
        cout << "========================================" << endl;
        cout << "PROCESSING COMMFLARES" << endl;
    };

    ProcessNodeCFAll(verbose);

    // transmit non-private info from Network to Nodes
    NetworkDataUpdateToNodes();

    // node timestamp summarizations 
    UpdateNodeData();
}


/*
// description 
iterates through network nodes and clears those that are bankrupt. 
*/ 
void Network::ClearNullNodes() {
   set<int> keep = {}; 

    
   for (auto x: contents) {
       if (x.second == nullptr) {
           delete x.second; 
           continue;
        }

       if ( (x.second)->GetWorth() > float(0.000000)) {
           keep.insert(x.first); 
       } else {
           x.second->ShutDown();
           delete contents[x.first];
           deceased.push_back(x.first);  
       }
   }

    map<int, NVMBNode*> newContents = map<int,NVMBNode*>(); 

    for (auto k: keep) {
        newContents[k] = contents[k]; 
    } 

    contents.clear();
    contents = map<int,NVMBNode*>(); 
    contents = newContents; 
}

/*
// description
runs each node's current task for one timestamp.

  *example*
- task could be
*** discovery or commerce, travel 1 unit.
*/
void Network::RunOneTimestamp(int verbose) {
    if (mortality) {
        ClearNullNodes(); 
    }

    if (verbose) {
        cout << "TTIMESTAMP " << timestamp << endl; 
    }

    // update max connectivity
    SetMaxConnectivity();

    /// TODO: concurrency here.
    map<int, NVMBNode*>::iterator it;
    NVMBNode* n;

    // iterate through nodes, clear their transmission rate, update their competitors
    /// TODO: refactor this chunk 
    NodeStrategos* ns;
    set<int> comp;
    for (auto x: contents) {
        (x.second)->ClearTransmissionAndResponse(); 
        comp = (x.second)->FindCompetitors();
        ns = (x.second)->GetStrategy();
        ns->UpdateCompetitors(comp);
    }

    // run contract activities 
    ActiveContractExecution(); 

    for (it = contents.begin(); it != contents.end(); it++) {
        n = it->second;
        n->ProcessOne(verbose);
    };

   

    // update
    UpdateNodes(verbose);
    WriteOutToFile(true); 
}

void Network::Run(int verbose, int numRounds) { 

    PrerunProcess(); 
    if (numRounds > 0) {
        while (IsAlive() && numRounds) {
            if (verbose == 1 || verbose == 2) {cout << "** TIMESTAMP:\t" << timestamp << endl;}; 
            RunOneTimestamp(verbose); 
            numRounds--; 
        }
    } else {
        while (IsAlive()) {
             if (verbose == 1 || verbose == 2) {cout << "** TIMESTAMP:\t" << timestamp << endl;}; 
            RunOneTimestamp(verbose); 
        }        
    }
    
    WriteOutToFile(false); 
    ShutDown();    
}

void Network::RunRandomNodeNatureVars(int verbose, int numRounds, float switchFrequency) {
    
    PrerunProcess(); 
    if (numRounds > 0) {
        while (IsAlive() && numRounds) {
            if (verbose == 1 || verbose == 2) {cout << "** TIMESTAMP:\t" << timestamp << endl;}; 
            RunOneTimestamp(verbose); 
            UpdateNodeNaturesAtRandom(switchFrequency);
            numRounds--; 
        }
    } else {
        while (IsAlive()) {
             if (verbose == 1 || verbose == 2) {cout << "** TIMESTAMP:\t" << timestamp << endl;}; 
            RunOneTimestamp(verbose); 
            UpdateNodeNaturesAtRandom(switchFrequency); 
        }        
    }
    
    WriteOutToFile(false); 
    ShutDown();  
} 

void Network::UpdateNodeNaturesAtRandom(float freq) {

    for (auto c: contents) {
        UpdateNodeNatureAtRandom(c.first, freq); 
    }
}

void Network::UpdateNodeNatureAtRandom(int nodeId, float freq) { 

    if (contents[nodeId] == nullptr) {
        return; 
    }

    float rf = RandomFloat();
    if (rf > freq) {
        return; 
    }

    map<string,string> nodeNature = MakeRandomNodeNatureVariables(); 

    auto strategy = contents[nodeId]->GetStrategy(); 
    strategy->competition = stof(nodeNature["competition"]); 
    strategy->greed = stof(nodeNature["greed"]); 
    strategy->negotiation = stof(nodeNature["negotiation"]); 
    strategy->negotiation = stof(nodeNature["growth"]); 
    strategy->CalculateNormedNatureVar(); 
}


void Network::PrerunProcess() {
    for (auto c: contents) {

        if (c.second == nullptr) {
            cout << "NODE " << c.first << " NULL" << endl; 
            continue;
        }

        auto npu = c.second->npu; 
        npu->LoadInitialPaths(); 
    }

    for (auto c: contents) {
        cout << "NODEEE " << c.first << endl; 
        (c.second->npu)->DisplayBestPaths(); 
    }
}


void Network::ShutDown() {
    for (auto c: contents) {
        if (c.second == nullptr) {
            continue; 
        }
        cout << "shutting down " << c.first << endl; 
        c.second->ShutDown(); 
    }
}

/// TODO: code this. 
void Network::Reopen() {

}

void Network::ActivitySummaryUpdate() {
    
    for (auto x: contents) {
        if (x.second == nullptr) {continue;} 
        auto pp = (x.second)->GetTransmission();
        (x.second)->LogTimestampAttribute(); 
  }
}

void Network::WriteOutToFile(bool threshold) {
    for (auto c: contents) {
        if (c.second == nullptr) {continue;}
        cout << c.first << " writing out to file" << endl; 
        c.second->WriteOutToFile(threshold);  
    }
}

bool Network::IsAlive() {
    
    int x = std::accumulate(contents.begin(), contents.end(),
        0.0,
        [](int t, pair<int, NVMBNode*> r)
        {
        return t + (r.second != nullptr) ? 1: 0;
        }); 

    return (Size() == 0) ? false : true; 
}

vector<int> Network::ActiveNodes() {

    vector<int> output = vector<int>(); 

    for (auto c: contents) {
        if (c.second != nullptr) {output.push_back(c.first);}
    }

    return output; 
}