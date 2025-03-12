// SPDX-License-Identifier: GPL-2.0-or-later
pragma solidity =0.7.6;
pragma abicoder v2;

import '@openzeppelin/contracts/token/ERC20/IERC20.sol';

import './libraries/LowGasSafeMath.sol';
import './interfaces/IV1SwapRouter.sol';
import './base/ImmutableState.sol';
import './base/PeripheryPaymentsWithFeeExtended.sol';
import './libraries/Constants.sol';
import './libraries/PulseXLibraryV1.sol';

/// @title PulseX V1 Swap Router
/// @notice Router for stateless execution of swaps against PulseX V1
abstract contract V1SwapRouter is IV1SwapRouter, ImmutableState, PeripheryPaymentsWithFeeExtended {
    using LowGasSafeMath for uint256;

    // supports fee-on-transfer tokens
    // requires the initial amount to have already been sent to the first pair
    function _swapV1(address[] memory path, address _to) private {
        for (uint256 i; i < path.length - 1; i++) {
            (address input, address output) = (path[i], path[i + 1]);
            (address token0, ) = PulseXLibraryV1.sortTokens(input, output);
            IPulseXPair pair = IPulseXPair(PulseXLibraryV1.pairFor(factoryV1, input, output));
            uint256 amountInput;
            uint256 amountOutput;
            // scope to avoid stack too deep errors
            {
                (uint256 reserve0, uint256 reserve1, ) = pair.getReserves();
                (uint256 reserveInput, uint256 reserveOutput) =
                    input == token0 ? (reserve0, reserve1) : (reserve1, reserve0);
                amountInput = IERC20(input).balanceOf(address(pair)).sub(reserveInput);
                amountOutput = PulseXLibraryV1.getAmountOut(amountInput, reserveInput, reserveOutput);
            }
            (uint256 amount0Out, uint256 amount1Out) =
                input == token0 ? (uint256(0), amountOutput) : (amountOutput, uint256(0));
            address to = i < path.length - 2 ? PulseXLibraryV1.pairFor(factoryV1, output, path[i + 2]) : _to;
            pair.swap(amount0Out, amount1Out, to, new bytes(0));
        }
    }

    /// @inheritdoc IV1SwapRouter
    function swapExactTokensForTokensV1(
        uint256 amountIn,
        uint256 amountOutMin,
        address[] calldata path,
        address to
    ) external payable override returns (uint256 amountOut) {
        // use amountIn == Constants.CONTRACT_BALANCE as a flag to swap the entire balance of the contract
        bool hasAlreadyPaid;
        if (amountIn == Constants.CONTRACT_BALANCE) {
            hasAlreadyPaid = true;
            amountIn = IERC20(path[0]).balanceOf(address(this));
        }

        pay(
            path[0],
            hasAlreadyPaid ? address(this) : msg.sender,
            PulseXLibraryV1.pairFor(factoryV1, path[0], path[1]),
            amountIn
        );

        // find and replace to addresses
        if (to == Constants.MSG_SENDER) to = msg.sender;
        else if (to == Constants.ADDRESS_THIS) to = address(this);

        uint256 balanceBefore = IERC20(path[path.length - 1]).balanceOf(to);

        _swapV1(path, to);

        amountOut = IERC20(path[path.length - 1]).balanceOf(to).sub(balanceBefore);
        require(amountOut >= amountOutMin, 'Too little received');
    }

    /// @inheritdoc IV1SwapRouter
    function swapTokensForExactTokensV1(
        uint256 amountOut,
        uint256 amountInMax,
        address[] calldata path,
        address to
    ) external payable override returns (uint256 amountIn) {
        amountIn = PulseXLibraryV1.getAmountsIn(factoryV1, amountOut, path)[0];
        require(amountIn <= amountInMax, 'Too much requested');

        pay(path[0], msg.sender, PulseXLibraryV1.pairFor(factoryV1, path[0], path[1]), amountIn);

        // find and replace to addresses
        if (to == Constants.MSG_SENDER) to = msg.sender;
        else if (to == Constants.ADDRESS_THIS) to = address(this);

        _swapV1(path, to);
    }
}
