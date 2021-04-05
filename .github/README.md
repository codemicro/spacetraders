# spacetraders

*An automatic client to play [SpaceTraders](https://spacetraders.io)*

![License](https://img.shields.io/github/license/codemicro/spacetraders) ![Lines of code](https://img.shields.io/tokei/lines/github/codemicro/spacetraders) [![Go Report Card](https://goreportcard.com/badge/github.com/codemicro/spacetraders)](https://goreportcard.com/report/github.com/codemicro/spacetraders)

---

### Not ready for use

### Strategy

* Take out loan (not yet automated)
* Purchase two ships (not yet automated)
* Pick the faster of the two ships (this bit not automated), visit everything that has a marketplace with no cargo
  * This is the scout ship
  * The scout ship makes a note of the marketplace on each planet visit
* Wait until the first ship has reached a second planet
* Run flights with the remaining ships for the best route with the best cargo from their current location (**WIP**)
  * Base this on `credits earned on sale / distance * volume` (credits per distance volume)
  * Ensure all planet marketplaces are re-indexed whenever arriving at a new planet
* When the scout has visited all planets with marketplaces, it then joins the trading fleet with the rest of the ships

### Not yet automated

* New account start
  * Taking out a loan, buying the first ships, etc
* Loan repayment
* Ship purchasing/scrapping
  * Some lower capacity ships will need to be scrapped as more ships are bought online due to ratelimit concerns