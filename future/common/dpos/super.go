package dpos

import (
	"github.com/Futuremine-chain/future/future/types"
	"github.com/Futuremine-chain/future/tools/arry"
)

// initialCandidates the first super node of the block generation cycle.
// The first half is the address of the block, the second half is the id of the block node
var genesisSuperList = types.Candidates{
	Members: []*types.Member{
		{
			Signer: arry.StringToAddress("FMizWwybDdxE9wWtbeiukj4ixcnMGLaT5La"),
			PeerId: "16Uiu2HAmTXRT6srXsTgy2kpzebMPPh5AqzvcLfFga74fcQ8BaEAH",
			Weight: 0,
		},
		{
			Signer: arry.StringToAddress("FMejc9bjiTeQzKQG9fSDPGdsRzzEdEQe6se"),
			PeerId: "16Uiu2HAkwKrbmaz3WRPjdJZbEBDCj412auZPoBCr3cpDViztzcX6",
			Weight: 0,
		},
		{
			Signer: arry.StringToAddress("FMhz3oBLYwRtRPXVAz869XFkYYGsEAhjiew"),
			PeerId: "16Uiu2HAmPQppQsHmVgJNqLPxTDNVy6ucuHt5i6eXMZ7Nn1Uz6Jwr",
			Weight: 0,
		},
		{
			Signer: arry.StringToAddress("FMcNASRUDr5symmsvQyWJ4iXFJYoEJLnPMz"),
			PeerId: "16Uiu2HAmGCas2irGWZNmSsqhCpNXXxapTEEtAsBhHENrCZB3nPzj",
			Weight: 0,
		},
		{
			Signer: arry.StringToAddress("FMdayTLVXkpDgHyB5kTzJEso92EuFHavbhp"),
			PeerId: "16Uiu2HAmFHNR3k254oHErF9FZXjCo1HZwHNqcP9vCkLgAL9CSwoP",
			Weight: 0,
		},
		{
			Signer: arry.StringToAddress("FMoTe9q8Duj6yQesQuTbgFuLk7hh19gbut9"),
			PeerId: "16Uiu2HAmUbzwyv33P135n3u7WAnfFCc7TJENKre53aXBPzYCJw4J",
			Weight: 0,
		},
		{
			Signer: arry.StringToAddress("FMZGtCTtYzpTZfhpJN2rwBFwoVVx8yctVkQ"),
			PeerId: "16Uiu2HAkwbnpzmXn2KSW5VZAaikKrWYUjYtGYx799cGj3Xg3jpgp",
			Weight: 0,
		},
		{
			Signer: arry.StringToAddress("FMegukTco2m1S9Y4ebXM9kVpQ6jqGGZBwWv"),
			PeerId: "16Uiu2HAmVBC9Fct91M3ffK71o1rkob9kzjsLKDjLESQruPDdbAMx",
			Weight: 0,
		},
		{
			Signer: arry.StringToAddress("FMnUqdciErY8UkxgQgVcVP2EDV9YTEiZE5g"),
			PeerId: "16Uiu2HAmPYAvkHqmvVhDH7poHGfBWPFKfwycQ4jssThqYJUZ3NSf",
			Weight: 0,
		},
	},
}
