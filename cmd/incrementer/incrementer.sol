contract Incrementer {
    address creator;
    uint number;

    function Incrementer() public 
    {
        creator = msg.sender; 
        number = 0;
    }

    function increment() 
    {
        number = number + 1;
    }
    
    function getNumber() constant returns (uint) 
    {
        return number;
    }
    
   
    function kill() 
    { 
        if (msg.sender == creator)
            suicide(creator);  
    }
    
}
