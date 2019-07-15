package main

import (
	"io"
	"log"
	"os/exec"
	"strconv"
	"bytes"
	"crypto/sha1"
	"math"
	"fmt"
	)

func main() {
	//usage is just running and then doing git push -f (-f needed since it thinks we're commmiting old files)
	hashmsg:=""
	//setup blob hash object for tree
	cmd := exec.Command("git", "hash-object", "-w", "main.go")
	stdin, err := cmd.StdinPipe()
	_, err = io.WriteString(stdin, hashmsg)
	stdin.Close()
	hash, err := cmd.CombinedOutput()
	
	//set up tree
	treeSetup:="100644 blob " + string(hash[:len(hash)-1]) + "\tmain.go"
	cmd = exec.Command("git", "mktree")
	stdin, err = cmd.StdinPipe()
	_, err = io.WriteString(stdin, treeSetup)
	stdin.Close()
	tree, err := cmd.CombinedOutput()
	
	var out []byte
	currentNum:=0
    commitMsg:= ""
	//set date and time so they don't change our message constantly
	cmd = exec.Command("export", "GIT_COMMITTER_DATE=2019-07-15T00:00:00+0000")
	cmd = exec.Command("export", "GIT_AUTHOR_DATE=2019-07-15T00:00:00+0000")
	
	
	
	
	cmd = exec.Command("git", "commit-tree", string(tree[:len(tree)-1]))
	stdin, _ = cmd.StdinPipe()
	io.WriteString(stdin, "")
	
	stdin.Close()
	
	out, _ = cmd.CombinedOutput()
	currentNum+=1
	
	
	cmd=exec.Command("git", "reset", "--hard", string(out[:len(out)-1]))
	out, err = cmd.CombinedOutput()
	if err != nil {
	log.Fatal(err)
	}
	cmd = exec.Command("git", "cat-file", "commit", "HEAD")
	outCat, _ := cmd.CombinedOutput()
	
	commitMsg = "commit "+strconv.Itoa(len(outCat)+int(math.Log10(float64(currentNum)))+1)+"\000" + string(outCat)+strconv.Itoa(currentNum)
	cmd = exec.Command("sha1sum")
	stdin, err = cmd.StdinPipe()
	_, err = io.WriteString(stdin, commitMsg)
	stdin.Close()
	out, err = cmd.CombinedOutput()

	commitMsg = "commit "+strconv.Itoa(len(outCat)+int(math.Log10(float64(currentNum)))+1)+"\000" + string(outCat) + strconv.Itoa(currentNum)
	s := commitMsg
	h := sha1.New()
    h.Write([]byte(s))
    bs := h.Sum(nil)
	
	goal:=make([]byte,3)
	/////////////////////////////////////////////////////////////
	for (bytes.Equal(bs[:3],goal)) == false {
	commitMsg = "commit "+strconv.Itoa(len(outCat)+int(math.Log10(float64(currentNum)))+1)+"\000" + string(outCat) + strconv.Itoa(currentNum)    
	s = commitMsg
	h = sha1.New()
    h.Write([]byte(s))
    bs = h.Sum(nil)
	
	currentNum+=1
	}
	//
	
	currentNum-=1
	cmd = exec.Command("git", "commit-tree", string(tree[:len(tree)-1]))
	
	stdin, err = cmd.StdinPipe()
	_, err = io.WriteString(stdin, strconv.Itoa(currentNum))
	stdin.Close()
	out, err = cmd.CombinedOutput()
	if err!=nil {
	log.Fatal(err)
	}
	
	fmt.Printf("%s/n",out)
	cmd = exec.Command("git", "reset", "--hard", string(out[:len(out)-2]))
	out, err = cmd.CombinedOutput()
	if err!=nil {
	log.Fatal(err)
	}
	}
