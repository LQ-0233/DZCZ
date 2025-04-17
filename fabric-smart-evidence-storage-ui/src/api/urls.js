const methodsTemplate = {
    "get": {
        method: "get"
    },
    "post": {
        method: "post"
    },
    "put": {
        method: "put"
    },
    "patch": {
        method: "patch"
    },
    "blob": {
        method: "get",
        responseType: "blob"
    }
};
   
const urlsTemplate = {

    "user": {
        "post": {
            "update": "user/update",
            "changePwd": "user/changePwd",
        },
        "get": {
            "list": "user/list",
        }
    },
    "common": {
        "post": {
            "login": "user/login",
            "register": "user/register",
            "upload": "common/upload",
        },
        "blob": {
            "download": "common/download"
        }
    },
    "evidence": {
        "post": {
            "create": "evidence/create",
    
            "authorize": "evidence/authorize",
            "cancelAuthorize": "evidence/cancelAuthorize",
            "view": "evidence/view",
        },
        "get": {
            "userEvidence": "evidence/userEvidence",
            "authorizedList": "evidence/authorizedList",
            "viewRecordList": "evidence/viewRecordList",
            "receivedAuthorizedList": "evidence/receivedAuthorizedList",
            "authorizeUserList": "evidence/authorizeUserList",
        }
    },
};

const urls = {};
Object.keys(urlsTemplate).forEach(function (group) {
    Object.keys(urlsTemplate[group]).forEach(function (methodName) {
        Object.keys(urlsTemplate[group][methodName]).forEach(function (name) {
            if (typeof urlsTemplate[group][methodName][name] === 'string') {
                urlsTemplate[group][methodName][name] = {
                    url: urlsTemplate[group][methodName][name]
                };
            }
            if (urls[group] === undefined) {
                urls[group] = {};
            }
            if (urls[group][methodName] === undefined) {
                urls[group][methodName] = {};
            }
            urls[group][methodName][name] = Object.assign({}, methodsTemplate[methodName], urlsTemplate[group][methodName][name]);
        });
    });
});
console.log(urls)
export default urls