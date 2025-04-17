'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');
const { v4: uuidv4 } = require('uuid');

class MyWorkload extends WorkloadModuleBase {
    constructor() {
        super();
    }

    async initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext) {
        await super.initializeWorkloadModule(workerIndex, totalWorkers, roundIndex, roundArguments, sutAdapter, sutContext);
       
    }

    async submitTransaction() {
        const uniqueId = uuidv4();

        console.log(`Worker ${this.workerIndex}: register user ${uniqueId}`);
        const request = {
                contractId: this.roundArguments.contractId,
                contractFunction: 'Register',
                invokerIdentity: 'User1',
                contractArguments: [`${uniqueId}`,"$2a$12$eJSlWpWr4mATJzX3n09NHOEpcpo34gt2KHrYNIV/FQVT4Pc/WOTyu","test","2"],
                readOnly: false
        };

         await this.sutAdapter.sendRequests(request);
    }


}

function createWorkloadModule() {
    return new MyWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;