pragma solidity ^0.4.11;
contract Survey {
    enum Answer{A,B,C,D}

    string public topic;
    
    bool finalized;
    address owner;
    uint[4] answerVotes;
    mapping (address => bool) votedUsers;

    function Survey(string topic_) public {
        owner = msg.sender; 
        topic = topic_;
    }

    function vote(Answer a) public {
        if (finalized) {
            revert();
        }

        if (votedUsers[msg.sender]) {
            revert();
        }

        votedUsers[msg.sender] = true;
        answerVotes[uint(a)] += 1;
    }
    
    function mostVoted() public constant returns (Answer) {
        Answer bestAnswer = Answer.A;
        uint votes = answerVotes[0];
        for (uint i = 1; i < 4; i++) {
            if (votes < answerVotes[i]) {
                votes = answerVotes[i];
                bestAnswer = Answer(i);
            }
        }
        return bestAnswer;
    }

    function voteForAnswer(Answer a) public constant returns (uint) {
        return answerVotes[uint(a)];
    }

    modifier onlyOwner() {
        require(msg.sender == owner);
        _;  
    }

    function finalize() onlyOwner public {
        finalized = true;
    }
}
