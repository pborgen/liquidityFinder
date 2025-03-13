# Liquidity Finder API Documentation

## Base URL

```
http://localhost:8080/api/v1
```

## Authentication

All API endpoints require an API key to be sent in the header:

```
X-API-Key: your-api-key
```

## Endpoints

### Pairs

#### GET /pairs

Get all trading pairs with optional filtering.

Query Parameters:

```
dex: string (optional) - Filter by DEX (e.g., "pulsex", "uniswap")
minLiquidity: number (optional) - Minimum liquidity in USD
hasValidTax: boolean (optional) - Filter for tax-validated tokens only
network: string (optional) - Network name (e.g., "pulsechain", "ethereum")
page: number (optional) - Page number for pagination
limit: number (optional) - Results per page (default: 20)
```

Response:

```json
{
  "success": true,
  "data": {
    "pairs": [
      {
        "id": "string",
        "dex": "string",
        "token0": {
          "address": "string",
          "symbol": "string",
          "name": "string",
          "decimals": "number"
        },
        "token1": {
          "address": "string",
          "symbol": "string",
          "name": "string",
          "decimals": "number"
        },
        "reserves": {
          "token0": "string",
          "token1": "string"
        },
        "liquidityUSD": "number",
        "hasTaxToken": "boolean",
        "lastUpdated": "string"
      }
    ],
    "pagination": {
      "currentPage": "number",
      "totalPages": "number",
      "totalResults": "number"
    }
  }
}
```

#### GET /pairs/{pairAddress}

Get detailed information about a specific trading pair.

Response:

```json
{
  "success": true,
  "data": {
    "pair": {
      "address": "string",
      "dex": "string",
      "token0": {
        "address": "string",
        "symbol": "string",
        "name": "string",
        "decimals": "number",
        "price": "number",
        "volume24h": "number"
      },
      "token1": {
        "address": "string",
        "symbol": "string",
        "name": "string",
        "decimals": "number",
        "price": "number",
        "volume24h": "number"
      },
      "reserves": {
        "token0": "string",
        "token1": "string"
      },
      "liquidityUSD": "number",
      "volume24h": "number",
      "fees24h": "number",
      "hasTaxToken": "boolean",
      "lastUpdated": "string"
    }
  }
}
```

### Tokens

#### GET /tokens

Get all tokens with optional filtering.

Query Parameters:

```
isTaxToken: boolean (optional) - Filter for tax tokens
minLiquidity: number (optional) - Minimum liquidity in USD
network: string (optional) - Network name
page: number (optional) - Page number
limit: number (optional) - Results per page
```

Response:

```json
{
  "success": true,
  "data": {
    "tokens": [
      {
        "address": "string",
        "name": "string",
        "symbol": "string",
        "decimals": "number",
        "totalLiquidity": "number",
        "price": "number",
        "volume24h": "number",
        "isTaxToken": "boolean",
        "taxPercentage": "number"
      }
    ],
    "pagination": {
      "currentPage": "number",
      "totalPages": "number",
      "totalResults": "number"
    }
  }
}
```

#### GET /tokens/{tokenAddress}

Get detailed information about a specific token.

Response:

```json
{
  "success": true,
  "data": {
    "token": {
      "address": "string",
      "name": "string",
      "symbol": "string",
      "decimals": "number",
      "totalLiquidity": "number",
      "price": "number",
      "volume24h": "number",
      "isTaxToken": "boolean",
      "taxPercentage": "number",
      "pairs": [
        {
          "address": "string",
          "dex": "string",
          "liquidityUSD": "number",
          "volume24h": "number"
        }
      ]
    }
  }
}
```

### Analytics

#### GET /analytics/overview

Get overview statistics across all DEXes.

Response:

```json
{
  "success": true,
  "data": {
    "totalPairs": "number",
    "totalTokens": "number",
    "totalLiquidityUSD": "number",
    "volume24hUSD": "number",
    "taxTokenPercentage": "number",
    "topDexes": [
      {
        "name": "string",
        "liquidityUSD": "number",
        "volume24h": "number",
        "pairs": "number"
      }
    ]
  }
}
```

#### GET /analytics/liquidity-changes

Get recent significant liquidity changes.

Query Parameters:

```
timeframe: string (optional) - Time window ("1h", "24h", "7d")
minChangePercent: number (optional) - Minimum change percentage
limit: number (optional) - Number of results
```

Response:

```json
{
  "success": true,
  "data": {
    "changes": [
      {
        "pair": {
          "address": "string",
          "dex": "string",
          "token0Symbol": "string",
          "token1Symbol": "string"
        },
        "liquidityChangePct": "number",
        "liquidityChangeUSD": "number",
        "timestamp": "string"
      }
    ]
  }
}
```

### Mempool

#### GET /mempool/transactions

Get recent mempool transactions (websocket recommended for real-time updates).

Query Parameters:

```
type: string (optional) - Transaction type ("swap", "addLiquidity", "removeLiquidity")
dex: string (optional) - DEX name
limit: number (optional) - Number of results
```

Response:

```json
{
  "success": true,
  "data": {
    "transactions": [
      {
        "hash": "string",
        "type": "string",
        "dex": "string",
        "pair": {
          "address": "string",
          "token0Symbol": "string",
          "token1Symbol": "string"
        },
        "value": "number",
        "timestamp": "string"
      }
    ]
  }
}
```

## WebSocket Endpoints

### ws://localhost:8080/ws/mempool

Real-time mempool transaction updates.

Message Format:

```json
{
  "type": "transaction",
  "data": {
    "hash": "string",
    "type": "string",
    "dex": "string",
    "pair": {
      "address": "string",
      "token0Symbol": "string",
      "token1Symbol": "string"
    },
    "value": "number",
    "timestamp": "string"
  }
}
```

### ws://localhost:8080/ws/pairs/{pairAddress}

Real-time updates for a specific pair.

Message Format:

```json
{
  "type": "update",
  "data": {
    "reserves": {
      "token0": "string",
      "token1": "string"
    },
    "price": "number",
    "liquidityUSD": "number",
    "timestamp": "string"
  }
}
```

## Error Responses

All endpoints return errors in the following format:

```json
{
  "success": false,
  "error": {
    "code": "string",
    "message": "string",
    "details": "object (optional)"
  }
}
```

Common Error Codes:

- `400` - Bad Request
- `401` - Unauthorized
- `403` - Forbidden
- `404` - Not Found
- `429` - Too Many Requests
- `500` - Internal Server Error
