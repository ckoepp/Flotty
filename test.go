package main


import (
  //"net"
  "github.com/nfnt/resize"
  "strconv"
  "io"
  "fmt"
  "os"
  "image"
  "image/color"
  "image/png"
  _ "image/jpeg"
  _ "image/gif"
)

func check(err error) {
  if err != nil {
    fmt.Println(err)
    panic(err)
  }
}

type PbmImage struct {
  width int
  height int
  content []bool
}

func NewPbmImage(width int, height int) *PbmImage {
  r := new(PbmImage)
  r.width = width
  r.height = height
  r.content = make([]bool, width * height)
  return r
}

func (pi *PbmImage) encode(writers []io.Writer) {

  width := pi.width / len(writers);

  // pbm header
  for _, w := range writers {
    w.Write([]byte("P1 "))
    w.Write([]byte(strconv.Itoa(width) + " "))
    w.Write([]byte(strconv.Itoa(pi.height)+ "\n"))
  }


  // pbm content
  w := writers[0]
  counter := 0
  wrr := 0

  for x := 0; x < len(pi.content); x++ {

    // counter
    if counter >= width - 1 {
      if wrr < len(writers) - 1 {
        wrr++
      } else {
        wrr = 0
      }
      w = writers[wrr]

      counter = 0
    } else {
      counter++
    }

    if pi.content[x] {
      w.Write([]byte("1 "))
    } else {
      w.Write([]byte("0 "))
    }
    if x % 70 == 0 {
      w.Write([]byte("\n"))
    }

  }
}

func (pi *PbmImage) store(x int, y int, black bool) {
  pi.content[y * pi.width + x] = black
}

func main() {
  filename := os.Args[1]
  original_image, err := os.Open(filename)
  check(err)
  defer original_image.Close()
  source_image,dataType, err := image.Decode(original_image)
  check(err)
  fmt.Printf("Opened %s which is of type %s\n", filename, dataType)

  source_image = resize.Resize(144, 120, source_image, resize.Lanczos3)
  bounds := source_image.Bounds()
  width, height := bounds.Max.X, bounds.Max.Y
  new_image := image.NewRGBA(image.Rectangle{image.Point{0,0}, image.Point{width, height}})

  // pbm header
  pbm := NewPbmImage(width, height)

  for x := 0; x < width; x++ {
    for y := 0; y < height; y++ {
      oldColor := source_image.At(x,y)
      r,g,b,alpha := oldColor.RGBA()

      // threshold calc
      avg := (r + g + b) / 3
      var c color.RGBA

      if uint8(alpha) == 0 || avg >= 22000 {
        pbm.store(x,y,false)
        c = color.RGBA{255,255,255,255}
      } else {
        pbm.store(x,y,true)
        c = color.RGBA{0,0,0, 255}
      }
      new_image.Set(x,y,c)
    }
  }

  myfile, _ := os.Create("test.png")
  //conn, err := net.Dial("udp6", "[2001:67c:20a1:1095:ba27:ebff:fe71:dd32]:2323")
  //check(err)
  png.Encode(myfile, new_image)

  pbmFile1, _ := os.Create("test1.pbm")
  pbmFile2, _ := os.Create("test2.pbm")
  pbmFile3, _ := os.Create("test3.pbm")
  pbm.encode([]io.Writer{pbmFile1, pbmFile2, pbmFile3})
  //pbm.encode(conn)
}
