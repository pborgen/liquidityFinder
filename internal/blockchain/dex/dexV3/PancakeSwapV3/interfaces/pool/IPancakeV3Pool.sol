// SPDX-License-Identifier: GPL-2.0-or-later
pragma solidity >=0.5.0;

import './pool/IPancakeV3PoolImmutables.sol';
import './pool/IPancakeV3PoolState.sol';
import './pool/IPancakeV3PoolDerivedState.sol';
import './pool/IPancakeV3PoolActions.sol';
import './pool/IPancakeV3PoolOwnerActions.sol';
import './pool/IPancakeV3PoolEvents.sol';

/// @title The interface for a 9mmSwap V3 Pool
/// @notice A 9mmSwap pool facilitates swapping and automated market making between any two assets that strictly conform
/// to the ERC20 specification
/// @dev The pool interface is broken up into many smaller pieces
interface IPancakeV3Pool is
    IPancakeV3PoolImmutables,
    IPancakeV3PoolState,
    IPancakeV3PoolDerivedState,
    IPancakeV3PoolActions,
    IPancakeV3PoolOwnerActions,
    IPancakeV3PoolEvents
{

}
