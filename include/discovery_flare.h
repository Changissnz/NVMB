// a flare is responsible for vertex discovery and civil increases
// The class <flare> interacts with <network> to determine which nodes to cross.
#ifndef FLARE_H
#define FLARE_H
#include <string>
#include <map>
#include <set>
#include <vector>
#include <sstream>
#include <iostream>
#include <iterator>
#include <utility>
#include <assert.h>
#include "extras.hpp"

class DiscoveryFlare
{

private:

  int vertexOwner;
  int vertexLocation;
  std::set<int> neighborIdentifiers;

  float worth;

  std::map<std::string, int> travelledEdges;
  std::vector<int> path; // path representation for simple undirected graph.

public:

  bool finished;

  DiscoveryFlare(int vo, int vl){
    vertexOwner = vo;
    vertexLocation = vl;
    neighborIdentifiers = std::set<int>();
    travelledEdges = std::map<std::string, int>();
    path = std::vector<int>();
    finished = false;
    worth = 0;
  };

  DiscoveryFlare(int vo) {
    vertexOwner = vo;
    vertexLocation = vo;
    neighborIdentifiers = std::set<int>();
    travelledEdges = std::map<std::string, int>();
    path = std::vector<int>();
    path.push_back(vo);
    finished = false;
    worth = 0;
  };

  int GetOwner();
  int GetLocation();
  float GetWorth();
  std::map<std::string, int> GetTravelledEdges() {return travelledEdges;};
  void UpdateLocation(int nl);

  // neighbor identifiers method
  std::set<int> GetNeighborIdentifiers();
  void UpdateNeighborIdentifers(std::set<int> ni);

  // methods for travelling
  void AddToPath(int i);
  std::vector<int> GetPath() {return path;};

  void DisplayPath();
  void DisplayNeighbors();
  void DisplayInfo();

  bool HasTravelledEdge(int v1, int v2);
  void MarkEdgeAsTravelled(std::string edge);

  DiscoveryFlare* MakeCopy();

  void UpdateWorth(float w);
};

#endif
