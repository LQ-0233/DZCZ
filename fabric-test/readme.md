
# 部署测试链码
```shell
# 在 test-network 目录下  

./network.sh deployCC  -c mychannel -ccn authentication -ccp ../authentication-center -ccl go

# 身份验证链码
./network.sh deployCC  -c mychannel -ccn evidence -ccp ../evidence-storage -ccl go
# 需初始化
./network.sh cc invoke -c mychannel -ccn evidence -ccic '{"Args":["InitLedger"]}'

# 调用测试
./network.sh cc query -c mychannel -ccn evidence -ccqc '{"Args":["VerifyBase64","AAAAAQAAAAAAAAABFCLoP4U81mNWDDfby5xF8dqyrlAKP30PfAUN8Q4ogw8=","xvmt5kAfl0eKGF0S95R+bCUNoFm9sfa4lkc/DuVxKXShoqZc361PSvFQ1aStWl+hckJxlcDpX5I+3dRVszgpTgbCBxNOCRJYRMkTQdQP/3Ow/KgO0wj9FcLIlCA000uV0eVGAL4XC/5U7QjToRo48EUnPqYYrXZSxjFRPEhWX2YAAAAAQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA="]}'

```

# 链存储TPS不低于xx，查询TPS不低于xx 30个节点下的交易时延不高于xxms 测试
```shell
# 测试的获取用户信息(查询)

npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networkConfig.yaml --caliper-benchconfig benchmark/getUserBenchmark.yaml --caliper-flow-only-test


# 测试 注册用户(存储/交易)

npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networkConfig.yaml --caliper-benchconfig benchmark/registerUserBenchmark.yaml --caliper-flow-only-test

# 身份验证

npx caliper launch manager --caliper-workspace ./ --caliper-networkconfig networkConfig.yaml --caliper-benchconfig benchmark/verifyBenchmark.yaml --caliper-flow-only-test
```

