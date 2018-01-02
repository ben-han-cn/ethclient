// This file is an automatically generated Go binding. Do not modify as any
// change will likely be lost upon the next re-generation!

package main

// SurveyABI is the input ABI used to generate the binding from.
const SurveyABI = "[{\"constant\":true,\"inputs\":[{\"name\":\"a\",\"type\":\"uint8\"}],\"name\":\"voteForAnswer\",\"outputs\":[{\"name\":\"\",\"type\":\"uint256\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[],\"name\":\"finalize\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"mostVoted\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"constant\":false,\"inputs\":[{\"name\":\"a\",\"type\":\"uint8\"}],\"name\":\"vote\",\"outputs\":[],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"constant\":true,\"inputs\":[],\"name\":\"topic\",\"outputs\":[{\"name\":\"\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"name\":\"topic_\",\"type\":\"string\"}],\"payable\":false,\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"

// SurveyBin is the compiled bytecode used for deploying new contracts.
const SurveyBin = `0x6060604052341561000f57600080fd5b6040516104b43803806104b4833981016040528080516001805461010060a860020a03191661010033600160a060020a031602179055919091019050600081805161005e929160200190610065565b5050610100565b828054600181600116156101000203166002900490600052602060002090601f016020900481019282601f106100a657805160ff19168380011785556100d3565b828001600101855582156100d3579182015b828111156100d35782518255916020019190600101906100b8565b506100df9291506100e3565b5090565b6100fd91905b808211156100df57600081556001016100e9565b90565b6103a58061010f6000396000f30060606040526004361061006c5763ffffffff7c0100000000000000000000000000000000000000000000000000000000600035041663224036aa81146100715780634bb278f31461009c578063781825d1146100b1578063b3f98adc146100e8578063bf63a57714610101575b600080fd5b341561007c57600080fd5b61008a60ff6004351661018b565b60405190815260200160405180910390f35b34156100a757600080fd5b6100af6101ad565b005b34156100bc57600080fd5b6100c46101e8565b604051808260038111156100d457fe5b60ff16815260200191505060405180910390f35b34156100f357600080fd5b6100af60ff60043516610242565b341561010c57600080fd5b6101146102db565b60405160208082528190810183818151815260200191508051906020019080838360005b83811015610150578082015183820152602001610138565b50505050905090810190601f16801561017d5780820380516001836020036101000a031916815260200191505b509250505060405180910390f35b6000600282600381111561019b57fe5b600481106101a557fe5b015492915050565b6001543373ffffffffffffffffffffffffffffffffffffffff90811661010090920416146101da57600080fd5b6001805460ff191681179055565b600254600090819060015b600481101561023a576002816004811061020957fe5b0154821015610232576002816004811061021f57fe5b0154915080600381111561022f57fe5b92505b6001016101f3565b509092915050565b60015460ff161561025257600080fd5b73ffffffffffffffffffffffffffffffffffffffff331660009081526006602052604090205460ff161561028557600080fd5b73ffffffffffffffffffffffffffffffffffffffff33166000908152600660205260409020805460ff1916600190811790915560028260038111156102c657fe5b600481106102d057fe5b018054909101905550565b60008054600181600116156101000203166002900480601f0160208091040260200160405190810160405280929190818152602001828054600181600116156101000203166002900480156103715780601f1061034657610100808354040283529160200191610371565b820191906000526020600020905b81548152906001019060200180831161035457829003601f168201915b5050505050815600a165627a7a7230582086703b362e7de857b9cac252a54895972e906feb718cb53003203b5b2d0f0aad0029`
