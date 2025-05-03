# Profit Calculator

Or loss calculator if you have a bad year. 

May 2025 Update: Still WIP. Can't use tool yet!

Profit Calculator is a niche project aimed to help users of Interactive Brokers LCC in Japan calculator their profit/loss for a given year. Finding our Profit/Loss is not trivial, as we need to find when we entered a trade in JPY terms, then find when we exited a trade in JPY terms. The implication of this is that it's possible to have had a profit despite a loss in USD terms, or vice versa.

As taxable events applies to equities, options, dividends, and interest, this tool will find the profit/loss per asset category.

## Prerequisites

- Interactive Brokers account
- Go 1.22 or higher

## How to Use

If you're from Wallstreet Bets, I have a Wendy's resume generator [here](https://youtu.be/dQw4w9WgXcQ?feature=shared).

- Compile program

```sh
go build .
```

- Run it

```sh
./profit-calc "input.csv"
```

NOTE: the input should be formatted in the default manner that Interactive Brokers outputs csvs. This can be obtained by going to Performance and Reports > Activity > By Year > Download as CSV.

## Input Format Notes

The various asset types in the outputted CSV have different numbers of columns. This script looks specifically for the following headers:

```
> Withholding Tax
[Withholding Tax Header Currency Date Description Amount Code]

> Interest
[Interest Header Currency Date Description Amount]

> Trades
[Trades Header DataDiscriminator Asset Category Currency Symbol Date/Time Quantity T. Price C. Price Proceeds Comm/Fee Basis Realized P/L MTM P/L Code]

> Dividends
[Dividends Header Currency Date Description Amount]
```
