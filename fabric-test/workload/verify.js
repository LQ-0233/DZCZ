'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');



class MyWorkload extends WorkloadModuleBase {
    constructor() {
        super();
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
        this.userCredentialBase64 = "AAAAAQAAAAAAAAABFCLoP4U81mNWDDfby5xF8dqyrlAKP30PfAUN8Q4ogw8="
        this.proofBase64 = "xvmt5kAfl0eKGF0S95R+bCUNoFm9sfa4lkc/DuVxKXShoqZc361PSvFQ1aStWl+hckJxlcDpX5I+3dRVszgpTgbCBxNOCRJYRMkTQdQP/3Ow/KgO0wj9FcLIlCA000uV0eVGAL4XC/5U7QjToRo48EUnPqYYrXZSxjFRPEhWX2YAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="
        

    }

    async submitTransaction() {
        const myArgs = {
            contractId: this.roundArguments.contractId,
            contractFunction: 'VerifyBase64',
            invokerIdentity: 'User1',
            contractArguments: [this.userCredentialBase64, this.proofBase64],
            readOnly: true
        };

        await this.sutAdapter.sendRequests(myArgs);
    }

}

function createWorkloadModule() {
    return new MyWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;