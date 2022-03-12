// SPDX-License-Identifier: MIT
pragma solidity ^0.8.10;
pragma experimental ABIEncoderV2; //Hint (or distraction): Allows returning arrays from functions

contract EscrowPayments {
    /*The space below is reserved for state variables (all made public) and the Constructor. For each
    state variable that you include please provide a comment specifying the part number, such as the one
    below. */
    //Part I


    struct item{
        string itemName;
        uint itemPrice;
        address itemOwner;
        bytes1 itemStatus;
    }

    address public owner;
    item[] public allItems;
    address public TTP;


    constructor() {
        owner = msg.sender;
        TTP = address(0);
    }
    

    function addItem(string calldata _itemName, uint _itemPrice) public {
        item memory newItem = item(_itemName, _itemPrice, address(0), 'A');
        allItems.push(newItem);
    }

    function listItems() public view returns (item[] memory) {
        return allItems;
    }

    function addTTP(address _TTP) public {
        require( TTP == address(0), "Trusted Third Party Already Assigned" );
        TTP = _TTP;
    }

    function buyItem(string calldata _itemName) public payable returns(string memory){
        for(uint it=0; it<allItems.length; it++ )
        {
            if(keccak256(bytes(_itemName)) == keccak256(bytes(allItems[it].itemName))){
                if( allItems[it].itemOwner != address(0))
                    return "Error: Item Already owned";
                if( msg.value <( allItems[it].itemPrice * (1 ether) ) )
                    return "Error: Not enough ether sent";
                if( allItems[it].itemStatus != 'A')
                    return "Error: Item not avalible";
                allItems[it].itemOwner = msg.sender;
                allItems[it].itemStatus = "P";
                string memory itemName = allItems[it].itemName;
                return string(abi.encodePacked( itemName, " now owned by: ", msg.sender )); 
            }
        }
        return "Item Does not Exist";
    }

    function confirmPurchase(string calldata _itemName, bool _status) public {
        for(uint it=0; it<allItems.length; it++ )
            if(keccak256(bytes(_itemName)) == keccak256(bytes(allItems[it].itemName)))
                if( allItems[it].itemOwner == msg.sender)
                {
                    if( _status == true )
                        allItems[it].itemStatus = "C";
                    else
                        allItems[it].itemStatus = "D";
                }
    }

    // You need to write a function named handleDispute, which takes item title and status as the
    // parameters, for this requirement. The function can only be invoked by the TTP and only if the status of the
    // item is Disputed.

    function handleDispute(string calldata _itemName, bytes1 _status ) public {
        if( msg.sender != TTP)
            return;

        for(uint it=0; it<allItems.length; it++ )
            if(keccak256(bytes(_itemName)) == keccak256(bytes(allItems[it].itemName)))
                if( allItems[it].itemStatus == "D")
                    allItems[it].itemStatus = _status;
    }


    // Once the item has been confirmed as a successful purchase, having status as Confirmed, the owner can
    // withdraw the price of the item from the contract. Similarly, once the item status has changed to Return, as
    // decided by the TTP for the disputed cases, the buyer can withdraw the price of the item from the contract.
    // You need to write a function named receivePayment, which takes item title as the parameter, for this
    // requirement. The function can only be invoked by owner and the buyer of the item. The status of the item
    // then changes to Expired.

    function receivePayment(string calldata _itemName) public {
        for(uint it=0; it<allItems.length; it++ ){
            if(keccak256(bytes(_itemName)) == keccak256(bytes(allItems[it].itemName))){
                // Case where the owner is receiving the payment
                if( msg.sender == owner )
                {
                    // check the status of the item
                    if( allItems[it].itemStatus == "C"){
                        // Send money to the owner
                        (payable(msg.sender)).transfer(allItems[it].itemPrice * (1 ether));
                        // Set the item to expired
                        allItems[it].itemStatus = "X";
                        return;
                    }
                }
                // Case where the buyer is receiving the payment
                if( msg.sender == allItems[it].itemOwner )
                {
                    // check the status of the item
                    if( allItems[it].itemStatus == "R"){
                        // Send money to the buyer
                        (payable(msg.sender)).transfer(allItems[it].itemPrice * (1 ether));
                        // Set the item to expired
                        allItems[it].itemStatus = "X";
                        return;
                    }
                }
            }
        }
    }
}