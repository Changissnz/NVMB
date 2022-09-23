#ifndef NODESTRATEGOS_HPP
#define NODESTRATEGOS_HPP
#include <map>
#include <vector>
#include <string>
#include "plan.hpp"

extern int FUTURE_PLAN_SIZE;
extern float ACTIVATION_THRESHOLD;

class NodeStrategos
{
private:

    float rank;
    std::string defaultPlan; // default plan is used in the case
    std::vector<Plan*> futureNodePlans;
    std::vector<Plan*> nodePlansCache; //
    std::vector<Plan*> plansInExecution; // where active discovery and commerce flares are
    std::vector<Plan*> finishedPlans;

public:
    int planIndex;
    int propIndex;
    int contractIndex;

    std::vector<Contract*> activeContracts;
    std::set<int> contractsConducted;

    // arguments for ContractReview
    static float averageNodeContractPathDistance;
    static float averageTransmissionMeasure;
    int estimatedDiameter; // max eccentricity of this node

    // information used for decision-making
    static std::map<std::string, float> edgeRisks; // key is edge, value is float
    static std::map<int, std::pair<int,float>> networkContractInfo;

    // node nature here
    float competition; // willingness to antagonize another node
    float greed; // willingness to do contracts with non-competitor nodes
    float normedC;
    float normedG;
    float negotiation;
    float growth;

    ///////// strings are labels for functions that instance uses
    // growth-performance, growth-relative, number of contracts, number of contracts relative, path score to max, path score to mean
    std::string competitionMeasure;
    // distance, transmission, distance+transmission, average distance, average transmission
    std::string contractMeasure;
    // best paths risk, predicted worth
    std::string newContractMeasure;
    // path scores-promise, skipped nodes-promise
    std::string bondAdvantageMeasure;
    // path scores,skipped nodes
    std::string bondDeletionAdvantageMeasure;
    std::set<int> competitors;

    /// information used for decision-making
    bool discoveryOn;

    NodeStrategos(float comp, float gr, float neg, float gro,
        std::string compMeasure, std::string contMeasure, std::string ncontMeasure,
        std::string bam, std::string bdam) {

        assert (comp >= 0 && comp <= 1);
        assert (gr >= 0 && gr <= 1);
        assert (neg >= 0 && neg <= 1);
        assert (gro >= 0 && gro <= 1);

        assert (compMeasure == "growth-performance" || compMeasure == "growth-relative"
        || compMeasure == "number of contracts" || compMeasure == "number of contracts relative"
        || compMeasure == "path score to max" || compMeasure == "path score to mean");

        assert (contMeasure == "distance" || contMeasure == "transmission" ||
        contMeasure == "distance+transmission" || contMeasure == "average distance" ||
        contMeasure == "average transmission");
        assert (ncontMeasure == "best paths risk" || ncontMeasure == "predicted worth");
        assert (bam == "path scores-promise" || bam == "skipped nodes-promise");
        assert (bdam == "skipped nodes" || bdam == "path scores");

        rank = 1.0;
        competition = comp;
        greed = gr;
        negotiation = neg;
        growth = gro;
        CalculateNormedNatureVar(); 

        competitionMeasure = compMeasure;
        contractMeasure = contMeasure;
        newContractMeasure = ncontMeasure;
        bondAdvantageMeasure = bam;
        bondDeletionAdvantageMeasure = bdam;

        planIndex = 0;
        propIndex = 0;
        contractIndex = 0;

        discoveryOn = false;
    }

    static NodeStrategos* MakeOne(std::map<std::string, std::string> userDefinedVariables = std::map<std::string, std::string>());

    void CalculateNormedNatureVar() {        
        float sum = competition + greed;

        normedC = ZeroDiv(competition, sum, 1.0); 
        normedG = ZeroDiv(greed, sum, 1.0);
    }

    void ClearFuturePlans() {
        futureNodePlans = std::vector<Plan*>();
    }

    void UpdateCompetitors(std::set<int> comp) {
        competitors = comp;
    }

    std::string GetDefaultPlan() {return defaultPlan;};

    std::vector<Plan*> GetFuturePlans() {
        return futureNodePlans; 
    }; 

    std::vector<Plan*> GetPlansInCache() {
        return nodePlansCache;
    };

    std::vector<Plan*> GetPlansInExecution() {
        return plansInExecution;
    };

    std::vector<Plan*> GetFinishedPlans() {
        return finishedPlans;
    };


    float GetRank() {return rank;};

    void SetRank(float r) {rank = r;};

    void SetCompetitiveness(float c) {
        if (c >= 0.0 && c <= 1.0) {
        competition = c;
        }
    };

    void SetGreediness(float g) {
        if (g >= 0.0 && g <= 1.0) {
        greed = g;
        }
    };

    void SetCompetitionMeasure(std::string competitionFunction) {
        assert (competitionFunction == "growth-perf" || competitionFunction == "growth-rel" ||
        competitionFunction == "contracts" || competitionFunction == "contracts-rel");
        competitionMeasure = competitionFunction;
    };

    void ChangeNodeNS(std::map<std::string, std::string> delta, std::string changeType); 
    void AddPlanToPriority(Plan* p);
    void AddPlans();
    void AddExecutedPlan(Plan* p);
    int FindPlanByIdentifier(std::string planStorage, int planId);
    void DeletePlan(std::string planType, int identifier);
    void FinishExecutedPlan(std::string indexType, int i);
    Plan* ChoosePlan();
    void InsertDefaultPlan(Plan* p); 
    void SwitchDiscoveryMode();

    Plan* MakeBondPlan(std::string bondType, int referenceNode, int sourceNode, int targetNode, float planScore);
    Plan* MakeContractPlan(int nodeId, int referenceNode, int otherNode, float planScore);
    Plan* BreakContractPlan(int nodeId, int referenceNode, int otherNode, float planScore, int breakPlanIdentifier); 
    Plan* MakeDiscoveryPlan(float planScore);
    void InsertFuturePlan(Plan* p, int verbose = 0);
    void DeleteActiveContract(int nodeId);

    void DisplayVariables();
    void DisplayAllPlans();
};

#endif
