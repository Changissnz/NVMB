#include "discovery_flare.h"
using namespace std;

int DiscoveryFlare::GetOwner() {
  return vertexOwner;
}

int DiscoveryFlare::GetLocation() {
  return vertexLocation;
}

void DiscoveryFlare::UpdateLocation(int nl) {
  vertexLocation = nl;
}

set<int> DiscoveryFlare::GetNeighborIdentifiers() {
    return neighborIdentifiers;
}

void DiscoveryFlare::UpdateNeighborIdentifers(set<int> ni) {
  neighborIdentifiers = ni;
}

void DiscoveryFlare::AddToPath(int i) {
  path.push_back(i);
}

void DiscoveryFlare::DisplayPath() {
  ostringstream vts;
  copy(path.begin(), path.end(), ostream_iterator<int>(vts, " "));
  cout << vts.str() << endl;
}

void DiscoveryFlare::DisplayNeighbors() {
  set<int>::iterator it;
  for (it = neighborIdentifiers.begin(); it != neighborIdentifiers.end(); it++) {
    cout << (*it);
  }
  cout << "\n";
}

bool DiscoveryFlare::HasTravelledEdge(int v1, int v2) {
  string s = NodePairToString(v1, v2);

  if (travelledEdges.find(s) != travelledEdges.end()) {
    return true;
  }
  return false;
}

void DiscoveryFlare::MarkEdgeAsTravelled(string edge) {
  travelledEdges[edge] = 1;
}

/*
// description
makes a copy of instance with empty neighbor identifiers
*/
DiscoveryFlare* DiscoveryFlare::MakeCopy() {
  DiscoveryFlare* df = new DiscoveryFlare(this->vertexOwner, this->vertexLocation);
  df->travelledEdges = this->travelledEdges; // copy-assignment
  df->path = this->path;
  df->UpdateWorth(worth);
  return df;
}

void DiscoveryFlare::UpdateWorth(float w) {
  worth += w;
}
