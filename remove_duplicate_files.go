package main

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "os"
    "path"
    "sync"
)


const ROOT_DIR string = "/Users/gleb/Desktop/brain_db_backup_2015_partial"


func main() {
    fmt.Printf("Running remove duplicates...\n")

    var wg sync.WaitGroup

    // Handle duplicate files in dir.
    wg.Add(1)
    removeDuplicatesInDir(ROOT_DIR, &wg)

    // If separated into monthly folders, recurse through.
    // This step is parallelizable. Just need to block at the end to
    // make sure all threads finish, using sync.WaitGroup.
    sortedDirContents, _ := ioutil.ReadDir(ROOT_DIR)
    for _, maybeDir := range sortedDirContents {
        if maybeDir.IsDir() {
            wg.Add(1)
            go removeDuplicatesInDir(path.Join(ROOT_DIR, maybeDir.Name()), &wg)
        }
    }

    // Wait for all to finish.
    wg.Wait()
}


func removeDuplicatesInDir(dirName string, wg *sync.WaitGroup) {
    // fmt.Println(dirName)
    sortedDirContents, _ := ioutil.ReadDir(dirName)

    lastFile := sortedDirContents[0]
    fullPath := path.Join(dirName, lastFile.Name())
    lastFileContents, fileReadError := ioutil.ReadFile(fullPath)
    if fileReadError != nil {
        fmt.Println(fileReadError)
        os.Exit(1)
    }
    numDuplicates := 0
    for _, newFile := range sortedDirContents {
        newFullPath := path.Join(dirName, newFile.Name())
        // fmt.Printf("Checking: %s\n", newFile.Name())
        if fullPath == newFullPath {
            continue
        }
        newFileContents, _ := ioutil.ReadFile(newFullPath)
        if bytes.Equal(lastFileContents, newFileContents) {
            fmt.Println("Found duplicate ...")
            fmt.Println(newFile.Name())
            numDuplicates++
            os.Remove(newFullPath)
        } else {
            lastFileContents = newFileContents
        }
    }
    fmt.Printf("%d files in %s\n", numDuplicates, dirName)
    wg.Done()
}
