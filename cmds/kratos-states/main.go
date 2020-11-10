package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type HeightVote struct {
	Round              string   `json:"round"`
	Prevotes           []string `json:"prevotes"`
	PrevotesBitArray   string   `json:"prevotes_bit_array"`
	Precommits         []string `json:"precommits"`
	PrecommitsBitArray string   `json:"precommits_bit_array"`
}

func (h HeightVote) String() string {
	res, _ := json.Marshal(h)
	return string(res)
}

type ConsensusState struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		RoundState struct {
			HeightRoundStep   string       `json:"height/round/step"`
			StartTime         time.Time    `json:"start_time"`
			ProposalBlockHash string       `json:"proposal_block_hash"`
			LockedBlockHash   string       `json:"locked_block_hash"`
			ValidBlockHash    string       `json:"valid_block_hash"`
			HeightVoteSet     []HeightVote `json:"height_vote_set"`
			Proposer          struct {
				Address string `json:"address"`
				Index   string `json:"index"`
			} `json:"proposer"`
		} `json:"round_state"`
	} `json:"result"`
}

type ConDumpStata struct {
	Jsonrpc string `json:"jsonrpc"`
	Result  struct {
		RoundState struct {
			Height     string `json:"height"`
			Round      string `json:"round"`
			Step       int    `json:"step"`
			Validators struct {
				Validators []struct {
					Address string `json:"address"`
					PubKey  struct {
						Type  string `json:"type"`
						Value string `json:"value"`
					} `json:"pub_key"`
					VotingPower      string `json:"voting_power"`
					ProposerPriority string `json:"proposer_priority"`
				} `json:"validators"`
				Proposer struct {
					Address string `json:"address"`
					PubKey  struct {
						Type  string `json:"type"`
						Value string `json:"value"`
					} `json:"pub_key"`
					VotingPower      string `json:"voting_power"`
					ProposerPriority string `json:"proposer_priority"`
				} `json:"proposer"`
			} `json:"validators"`
			Proposal struct {
				Type     int    `json:"Type"`
				Height   string `json:"height"`
				Round    string `json:"round"`
				PolRound string `json:"pol_round"`
				BlockID  struct {
					Hash  string `json:"hash"`
					Parts struct {
						Total string `json:"total"`
						Hash  string `json:"hash"`
					} `json:"parts"`
				} `json:"block_id"`
				Signature string `json:"signature"`
			} `json:"proposal"`
			ProposalBlock struct {
				Header struct {
					Version struct {
						Block string `json:"block"`
						App   string `json:"app"`
					} `json:"version"`
					ChainID     string `json:"chain_id"`
					Height      string `json:"height"`
					LastBlockID struct {
						Hash  string `json:"hash"`
						Parts struct {
							Total string `json:"total"`
							Hash  string `json:"hash"`
						} `json:"parts"`
					} `json:"last_block_id"`
					LastCommitHash     string `json:"last_commit_hash"`
					DataHash           string `json:"data_hash"`
					ValidatorsHash     string `json:"validators_hash"`
					NextValidatorsHash string `json:"next_validators_hash"`
					ConsensusHash      string `json:"consensus_hash"`
					AppHash            string `json:"app_hash"`
					LastResultsHash    string `json:"last_results_hash"`
					EvidenceHash       string `json:"evidence_hash"`
					ProposerAddress    string `json:"proposer_address"`
				} `json:"header"`
				Data struct {
					Txs []string `json:"txs"`
				} `json:"data"`
			} `json:"proposal_block"`
			ProposalBlockParts struct {
				CountTotal    string `json:"count/total"`
				PartsBitArray string `json:"parts_bit_array"`
			} `json:"proposal_block_parts"`
			LockedRound               string `json:"locked_round"`
			ValidRound                string `json:"valid_round"`
			CommitRound               string `json:"commit_round"`
			TriggeredTimeoutPrecommit bool   `json:"triggered_timeout_precommit"`
		} `json:"round_state"`
		Peers []struct {
			NodeAddress string `json:"node_address"`
			PeerState   struct {
				RoundState struct {
					Height                   string    `json:"height"`
					Round                    string    `json:"round"`
					Step                     int       `json:"step"`
					StartTime                time.Time `json:"start_time"`
					Proposal                 bool      `json:"proposal"`
					ProposalBlockPartsHeader struct {
						Total string `json:"total"`
						Hash  string `json:"hash"`
					} `json:"proposal_block_parts_header"`
					ProposalBlockParts interface{} `json:"proposal_block_parts"`
					ProposalPolRound   string      `json:"proposal_pol_round"`
					ProposalPol        string      `json:"proposal_pol"`
					Prevotes           string      `json:"prevotes"`
					Precommits         string      `json:"precommits"`
					LastCommitRound    string      `json:"last_commit_round"`
					LastCommit         string      `json:"last_commit"`
					CatchupCommitRound string      `json:"catchup_commit_round"`
					CatchupCommit      string      `json:"catchup_commit"`
				} `json:"round_state"`
				Stats struct {
					Votes      string `json:"votes"`
					BlockParts string `json:"block_parts"`
				} `json:"stats"`
			} `json:"peer_state"`
		} `json:"peers"`
	} `json:"result"`
}

func (c *ConDumpStata) ShortString(addr string) string {
	return fmt.Sprintf("\tAddr: %s,\th: %s,\tproposer: %s,\tround: %s\tappHash: %s",
		addr,
		c.Result.RoundState.Height,
		c.Result.RoundState.Proposal.BlockID.Hash,
		c.Result.RoundState.Round,
		c.Result.RoundState.ProposalBlock.Header.AppHash,
	)
}

func (c *ConsensusState) ShortString(addr string) string {
	var vset HeightVote
	if len(c.Result.RoundState.HeightVoteSet) > 0 {
		vset = c.Result.RoundState.HeightVoteSet[len(c.Result.RoundState.HeightVoteSet)-1]
	}

	return fmt.Sprintf("\tAddr: %s,\th: %s,\tlast: %s\tproposer: %s,\tvoteSet: %s",
		addr,
		c.Result.RoundState.HeightRoundStep,
		c.Result.RoundState.ProposalBlockHash,
		c.Result.RoundState.Proposer.Address,
		vset.Round)
}

func QueryToJSON(res interface{}, format string, args ...interface{}) error {
	path := fmt.Sprintf(format, args...)

	client := http.Client{}
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Encoding", "gzip")
	req.Header.Add("Accept-Encoding", "gzip,deflate")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Upgrade-Insecure-Requests", "1")

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "error by get with %s", path)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "error by read all with %s", path)
	}

	if resp.StatusCode == 200 {
		if err := json.Unmarshal(body, res); err != nil {
			return errors.Wrapf(err, "error by unmarshal with %s", path)
		}

		return nil
	} else {
		return fmt.Errorf("resp code by %d with %s", resp.StatusCode, path)
	}
}

type AddrResState struct {
	Address [2]string `json:"address"`
	Err     error
	Cs      ConsensusState `json:"cs"`
	CsDump  ConDumpStata   `json:"cd"`
}

func (r AddrResState) String() string {
	if r.Err == nil {
		return fmt.Sprintf("stat: %s\t::\t%s", r.Address[0], r.CsDump.ShortString(r.Address[1]))
	} else {
		return fmt.Sprintf("stat: %s\t::\tERROR by %s", r.Address[0], r.Err.Error())
	}
}

func reqNodeStat(address [2]string) AddrResState {
	res := AddrResState{
		Address: address,
	}

	//if err := QueryToJSON(&res.Cs, address[1]+"consensus_state"); err != nil {
	//		res.Err = errors.Wrapf(err, "query consensus_state")
	//		return res
	//	}

	if err := QueryToJSON(&res.CsDump, address[1]+"dump_consensus_state"); err != nil {
		res.Err = errors.Wrapf(err, "query dump_consensus_state")
		return res
	}

	return res
}

func reqAllNodes(address [][2]string) []AddrResState {
	res := make([]AddrResState, 0, len(address))
	var mutex sync.Mutex

	var wg sync.WaitGroup

	for _, addr := range address {
		wg.Add(1)
		go func(addrParam [2]string) {
			defer wg.Done()

			statRes := reqNodeStat(addrParam)

			mutex.Lock()
			res = append(res, statRes)
			mutex.Unlock()
		}(addr)
	}

	wg.Wait()

	return res
}

func main() {
	// TODO: use cmd
	addresses := [][2]string{
		{"name", "http://127.0.0.1:26657/"},
	}

	for {
		time.Sleep(2 * time.Second)
		res := reqAllNodes(addresses)

		fmt.Printf("Res in %s\n", time.Now())
		for _, re := range res {
			fmt.Printf("%s\n", re.String())
		}
		fmt.Printf("Res end %s\n\n\n\n", time.Now())
	}
}
