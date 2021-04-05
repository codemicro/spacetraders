# spacetraders

*An automatic client to play [SpaceTraders](https://spacetraders.io)*

![License](https://img.shields.io/github/license/codemicro/spacetraders) ![Lines of code](https://img.shields.io/tokei/lines/github/codemicro/spacetraders) [![Go Report Card](https://goreportcard.com/badge/github.com/codemicro/spacetraders)](https://goreportcard.com/report/github.com/codemicro/spacetraders)

---

### Not ready for use

### Strategy

* Take out starter loan (not yet automated)
* Purchase three ships (not yet automated)
* Designate two ships as probe ships
  * Send these two ships to two locations and sell all remaining fuel
  * Use these ships to enable access to the market information for that location
  * Collect market data every minute (or two?)
* Use data from probes to find the best value route for trading ships to fly
  * Base this on `credits earned on sale / distance * volume` (credits per distance volume)
  * Only fly to locations with known markets

### Not yet automated

* New account start
  * Taking out a loan, buying the first ships, etc
* Loan repayment
* Ship purchasing/scrapping
  * Some lower capacity ships will need to be scrapped as more ships are bought online due to ratelimit concerns
  
### `// TODO: ...`

* Take available funds into account when planning what cargo to purchase.
