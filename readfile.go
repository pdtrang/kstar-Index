package kstar

import (
     "os"
    "fmt"
    "bufio"
    "bytes"
    //"log"                                                                                                                                                                                                                                  
    "strings"
    "strconv"
    "sort"
//      "runtime"                                                                                                                                                                                                                            
//      "log"                                                                                                                                                                                                                                
)

type SNP struct{
    profile []string
}

func (f SNP) GetString() []string {
    return f.profile
}

func SaveIndex2(i_qgram [][]int, prefix []string, K int, path string){

//      memstats := new(runtime.MemStats)                                                                                                                                                                                                    
//      runtime.ReadMemStats(memstats)                                                                                                                                                                                                       
//      log.Printf("Save Index 2: Alloc %d\tTotalAlloc %d\tSys %d\tHeapAlloc %d\tHeapSys %d", memstats.Alloc, memstats.TotalAlloc, memstats.Sys, memstats.HeapAlloc, memstats.HeapSys)                                                       
        file, err := os.Create(path)
        if err != nil {
                fmt.Println("error created file")
        }
        defer file.Close()

        w := bufio.NewWriter(file)

        fmt.Fprintln(w, prefix[0], ", K =", K)
        for j := 0; j < len(i_qgram); j++ {
                if(len(i_qgram[j])>0){
                        fmt.Fprint(w, j, "\t")
                        for k := 0; k < len(i_qgram[j]); k++ {
                                fmt.Fprint(w, i_qgram[j][k], " ")
                        }
                        fmt.Fprint(w, "\n")
                }

        }

//      runtime.ReadMemStats(memstats)                                                                                                                                                                                                       
//      log.Printf("Save Index 2: Alloc %d\tTotalAlloc %d\tSys %d\tHeapAlloc %d\tHeapSys %d", memstats.Alloc, memstats.TotalAlloc, memstats.Sys, memstats.HeapAlloc, memstats.HeapSys)                                                       

        w.Flush()

}

func SaveIndex(i_qgram [][]int, prefix []string, K int, path string){
//      memstats := new(runtime.MemStats)                                                                                                                                                                                                    
//      runtime.ReadMemStats(memstats)                                                                                                                                                                                                       
//      log.Printf("Save Index 1: Alloc %d\tTotalAlloc %d\tSys %d\tHeapAlloc %d\tHeapSys %d", memstats.Alloc, memstats.TotalAlloc, memstats.Sys, memstats.HeapAlloc, memstats.HeapSys)                                                       

        file, err := os.Create(path)
        if err != nil {
                fmt.Println("error created file")
        }
        defer file.Close()

        w := bufio.NewWriter(file)

        fmt.Fprintln(w, prefix[0], ", K =", K)
        for j := 0; j < len(i_qgram); j++ {
                for k := 0; k < len(i_qgram[j]); k++ {
                        fmt.Fprint(w, i_qgram[j][k], " ")
                }
                fmt.Fprint(w, "\n")
        }

//      runtime.ReadMemStats(memstats)                                                                                                                                                                                                       
//      log.Printf("Save Index 1: Alloc %d\tTotalAlloc %d\tSys %d\tHeapAlloc %d\tHeapSys %d", memstats.Alloc, memstats.TotalAlloc, memstats.Sys, memstats.HeapAlloc, memstats.HeapSys)                                                       


        w.Flush()

}

func ReadLines(path string) ([]string, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }
  return lines, scanner.Err()
}

/* Read real VCF file */
func ReadVCF(sequence_file string) map[int]SNP {
    array := make(map[int]SNP)
    f,err := os.Open(sequence_file)
    if err != nil{
        fmt.Printf("%v\n",err)
        os.Exit(1)
    }

    defer f.Close()
    br := bufio.NewReader(f)
    //byte_array := bytes.Buffer{}                                                                                                                                                                                                           

    for{
        line , err := br.ReadString('\n')
        if err != nil {
            //fmt.Printf("%v\n",err)                                                                                                                                                                                                         
            break
        }
        if line[0]==byte('#') {
            //fmt.Printf("%s \n",line)                                                                                                                                                                                                       
        } else {
            sline := string(line)
            split := strings.Split(sline, "\t");
            //fmt.Printf("%s %s %s\n", split[1], split[3], split[4])                                                                                                                                                                         
            pos, _ := strconv.ParseInt(split[1], 10, 64)
            pos = pos - 1
            if len(split[4])>1 {
                alt := strings.Split(split[4], ",")
                t := make([]string, len(alt)+1)
                t[0] = split[3]
                for i:=0; i<len(alt); i++ {
                    if alt[i] == "<DEL>" {
                        t[i+1] = "."
                    } else {
                        t[i+1] = alt[i]
                    }
                }
                //sort.Strings(t)                                                                                                                                                                                                            
                //array[int(pos)] = SNP{t} // asign SNP at pos                                                                                                                                                                               
                tmp, ok := array[int(pos)]
                if ok {
                    t = append(t[:0], t[1:]...)
                    tmp.profile = append(tmp.profile, t...)
                } else {
                    tmp.profile = append(tmp.profile, t...)
                }
                sort.Strings(tmp.profile)
                array[int(pos)] = tmp // append SNP at pos                                                                                                                                                                                   
                //fmt.Printf("pos=%d %q \n", pos, alt)                                                                                                                                                                                       
            } else {
                //array[int(pos)] = SNP{[]string{split[3], split[4]}} // asign SNP at pos                                                                                                                                                    
                tmp, ok := array[int(pos)]
                if ok {
                    if split[4] == "<DEL>" {
                        tmp.profile = append(tmp.profile, ".")
                    } else {
                        tmp.profile = append(tmp.profile, split[4])
                    }
                } else {
                    if split[4] == "<DEL>" {
                        tmp.profile = append(tmp.profile, []string{split[3], "."}...)
                    } else {
                        tmp.profile = append(tmp.profile, []string{split[3], split[4]}...)
                    }
                }
                sort.Strings(tmp.profile)
                array[int(pos)]= tmp // append SNP at pos                                                                                                                                                                                    
                //fmt.Println(pos)                                                                                                                                                                                                           
            }
        }
    }
    return array
}

func ReadFASTA(sequence_file string) []byte {
    f,err := os.Open(sequence_file)
    if err != nil{
        fmt.Printf("%v\n",err)
        os.Exit(1)
    }

    defer f.Close()
    br := bufio.NewReader(f)
    byte_array := bytes.Buffer{}

    //line , err := br.ReadString('\n')                                                                                                                                                                                                      
        _ , isPrefix, err := br.ReadLine()
        if err != nil || isPrefix{
                fmt.Printf("%v\n",err)
                os.Exit(1)
        }
    //fmt.Printf("%s",line)                                                                                                                                                                                                                  

    for {
        line , isPrefix, err := br.ReadLine()
        if err != nil || isPrefix{
            break
        } else {
            byte_array.Write([]byte(line))
        }
    }
    //byte_array.Write([]byte("$"))                                                                                                                                                                                                          
    input := []byte(byte_array.String())
    return input
}
