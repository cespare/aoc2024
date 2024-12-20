package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"iter"
	"log"
	"net/http"
	"os"
	"runtime"
	"runtime/pprof"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/kr/pretty"
	"golang.org/x/exp/constraints"
)

const year = 2024

var solutions = make(map[int][]func(*problemContext))

func addSolutions(problem int, fns ...func(*problemContext)) {
	solutions[problem] = append(solutions[problem], fns...)
}

func findSolution(problem, solNumber int) (func(*problemContext), error) {
	solns, ok := solutions[problem]
	if !ok {
		return nil, fmt.Errorf("no solutions for problem %d", problem)
	}
	if solNumber > len(solns) {
		return nil, fmt.Errorf("problem %d only has %d solution(s)", problem, len(solns))
	}
	return solns[solNumber-1], nil
}

func parseProblem(name string) (problem, solNumber int, err error) {
	parts := strings.SplitN(name, ".", 2)
	solNumber = 1
	if len(parts) == 2 {
		var err error
		solNumber, err = strconv.Atoi(parts[1])
		if err != nil {
			return 0, 0, err
		}
	}
	problem, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, err
	}
	return problem, solNumber, nil
}

func main() {
	log.SetFlags(0)

	cpuProfile := flag.String("cpuprofile", "", "write CPU profile to `file`")
	printTimings := flag.Bool("t", false, "print timings")
	readStdin := flag.Bool("stdin", false, "read from stdin instead of default file")
	downloadInput := flag.Bool("dlinput", false, "download the input and store as <day>.txt")
	flag.Parse()

	if *printTimings && *cpuProfile != "" {
		log.Fatal("-t and -cpuprofile are incompatible")
	}
	if flag.NArg() != 1 {
		log.Fatalf("Usage: %s [flags] problem", os.Args[0])
	}
	if *downloadInput {
		downloadProblemInput(flag.Arg(0))
		return
	}
	problem, solNumber, err := parseProblem(flag.Arg(0))
	if err != nil {
		log.Fatalln("Bad problem number:", err)
	}
	fn, err := findSolution(problem, solNumber)
	if err != nil {
		log.Fatal(err)
	}
	ctx, err := newProblemContext(problem, *readStdin)
	if err != nil {
		log.Fatal(err)
	}
	defer ctx.close()

	if *cpuProfile != "" {
		f, err := os.Create(*cpuProfile)
		if err != nil {
			log.Fatal(err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Fatalln("Error writing CPU profile:", err)
			}
		}()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatalln("Error starting CPU profile:", err)
		}
		defer pprof.StopCPUProfile()

		ctx.profiling = true
		fn(ctx)
		return
	}

	ctx.timings.start = time.Now()
	fn(ctx)
	ctx.timings.done = time.Now()
	if *printTimings {
		ctx.printTimings()
	}
}

func downloadProblemInput(day string) {
	n, err := strconv.Atoi(day)
	if err != nil {
		log.Fatalf("Bad problem number %q", day)
	}
	if n < 0 || n > 25 {
		log.Fatalf("Day number %d out of range", n)
	}
	day = fmt.Sprintf("%02d", n)

	cookieText, err := os.ReadFile("sessioncookie.txt")
	if err != nil {
		log.Fatalln("Error reading session cookie:", err)
	}
	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, n)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	req.AddCookie(&http.Cookie{
		Name:  "session",
		Value: strings.TrimSpace(string(cookieText)),
	})
	req.Header.Set(
		"User-Agent",
		fmt.Sprintf("github.com/cespare/aoc%d by cespare@gmail.com", year),
	)

	f, err := os.OpenFile(day+".txt", os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
	if err != nil {
		log.Fatalln("Error opening input file:", err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalln("Error downloading input:", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("Non-200 status response (%d)", resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		if err == nil && len(body) > 0 {
			log.Printf("Response:\n%s", body)
		}
		os.Exit(1)
	}
	if _, err := io.Copy(f, resp.Body); err != nil {
		log.Fatalln("Error writing input file:", err)
	}

	if err := f.Close(); err != nil {
		log.Fatalln("Error writing input file:", err)
	}
}

type problemContext struct {
	f            *os.File
	needClose    bool
	profiling    bool
	profileStart time.Time
	l            *log.Logger

	timings struct {
		start time.Time
		load  time.Time
		part1 time.Time
		part2 time.Time
		done  time.Time
	}
}

func (ctx *problemContext) reportLoad() { ctx.timings.load = time.Now() }

func (ctx *problemContext) reportPart1(v ...interface{}) {
	ctx.timings.part1 = time.Now()
	args := append([]interface{}{"Part 1:"}, v...)
	ctx.l.Println(args...)
}

func (ctx *problemContext) reportPart2(v ...interface{}) {
	ctx.timings.part2 = time.Now()
	args := append([]interface{}{"Part 2:"}, v...)
	ctx.l.Println(args...)
}

func (ctx *problemContext) printTimings() {
	ctx.l.Println("Total:", ctx.timings.done.Sub(ctx.timings.start))
	t := ctx.timings.start
	if !ctx.timings.load.IsZero() {
		ctx.l.Println("  Load:", ctx.timings.load.Sub(t))
		t = ctx.timings.load
	}
	if !ctx.timings.part1.IsZero() {
		ctx.l.Println("  Part 1:", ctx.timings.part1.Sub(t))
		t = ctx.timings.part1
	}
	if !ctx.timings.part2.IsZero() {
		ctx.l.Println("  Part 2:", ctx.timings.part2.Sub(t))
	}
}

func newProblemContext(n int, readStdin bool) (*problemContext, error) {
	ctx := &problemContext{
		l: log.New(os.Stdout, "", 0),
	}
	if readStdin {
		ctx.f = os.Stdin
	} else {
		name := fmt.Sprintf("%02d.txt", n)
		f, err := os.Open(name)
		if err != nil {
			return nil, err
		}
		ctx.f = f
		ctx.needClose = true
	}
	return ctx, nil
}

func (ctx *problemContext) close() {
	if ctx.needClose {
		ctx.f.Close()
	}
}

func (ctx *problemContext) loopForProfile() bool {
	if ctx.profileStart.IsZero() {
		ctx.profileStart = time.Now()
		return true
	}
	if !ctx.profiling {
		return false
	}
	return time.Since(ctx.profileStart) < 5*time.Second
}

func (ctx *problemContext) scanner() scanner {
	return newScanner(ctx.f)
}

func (ctx *problemContext) int64s() []int64 {
	var ns []int64
	s := ctx.scanner()
	for s.scan() {
		ns = append(ns, s.int64())
	}
	return ns
}

func (ctx *problemContext) lines() []string {
	var lines []string
	s := ctx.scanner()
	for s.scan() {
		lines = append(lines, s.text())
	}
	return lines
}

func (ctx *problemContext) readAll() []byte {
	b, err := io.ReadAll(ctx.f)
	if err != nil {
		log.Fatalln("Read error:", err)
	}
	return b
}

func scanSlice[E any](ctx *problemContext, parse func(string) E) []E {
	var vs []E
	scanner := ctx.scanner()
	for scanner.scan() {
		vs = append(vs, parse(scanner.text()))
	}
	return vs
}

type scanner struct {
	s *bufio.Scanner
}

func newScanner(r io.Reader) scanner {
	return scanner{bufio.NewScanner(r)}
}

func (s scanner) scan() bool {
	if !s.s.Scan() {
		if err := s.s.Err(); err != nil {
			log.Fatalln("Scan error:", err)
		}
		return false
	}
	return true
}

func (s scanner) text() string {
	return s.s.Text()
}

func (s scanner) int64() int64 {
	n, err := strconv.ParseInt(s.text(), 10, 64)
	if err != nil {
		log.Fatalf("Bad int64 %q: %s", s.text(), err)
	}
	return n
}

func (s scanner) fields() []string {
	return strings.Fields(s.text())
}

type number interface {
	constraints.Integer | constraints.Float
}

func abs[T number](n T) T {
	if n < 0 {
		return -n
	}
	return n
}

type vec2 struct {
	x int64
	y int64
}

func (v vec2) String() string {
	return fmt.Sprintf("%d,%d", v.x, v.y)
}

func (v vec2) add(v1 vec2) vec2 {
	return vec2{v.x + v1.x, v.y + v1.y}
}

func (v vec2) sub(v1 vec2) vec2 {
	return vec2{v.x - v1.x, v.y - v1.y}
}

func (v vec2) mul(m int64) vec2 {
	return vec2{v.x * m, v.y * m}
}

func (v vec2) eltMul(v1 vec2) vec2 {
	return vec2{v.x * v1.x, v.y * v1.y}
}

func (v vec2) mag() int64 {
	return abs(v.x) + abs(v.y)
}

func (v vec2) inv() vec2 {
	return vec2{-v.x, -v.y}
}

var (
	north = vec2{0, -1}
	east  = vec2{1, 0}
	south = vec2{0, 1}
	west  = vec2{-1, 0}
)

var nesw = []vec2{north, east, south, west}

var box8 = []vec2{
	{-1, -1},
	{0, -1},
	{1, -1},
	{-1, 0},
	{1, 0},
	{-1, 1},
	{0, 1},
	{1, 1},
}

func (v vec2) neighbors4() []vec2 {
	neighbors := make([]vec2, 4)
	for i, d := range nesw {
		neighbors[i] = v.add(d)
	}
	return neighbors
}

func (v vec2) neighbors8() []vec2 {
	neighbors := make([]vec2, 8)
	for i, d := range box8 {
		neighbors[i] = v.add(d)
	}
	return neighbors
}

type mat2 struct {
	a00, a01 int64
	a10, a11 int64
}

func (v vec2) matMul(m mat2) vec2 {
	return vec2{
		v.x*m.a00 + v.y*m.a01,
		v.x*m.a10 + v.y*m.a11,
	}
}

var (
	turnCW = mat2{
		0, -1,
		1, 0,
	}
	turnCCW = mat2{
		0, 1,
		-1, 0,
	}
	turn180 = mat2{
		-1, 0,
		0, -1,
	}
)

func parseInt(s string) int64 {
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return n
}

func parseUint(s string) uint64 {
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		panic(err)
	}
	return n
}

func sscanf(s, format string, args ...any) {
	if _, err := fmt.Sscanf(s, format, args...); err != nil {
		panic(fmt.Sprintf("sscanf error: %s", err))
	}
}

type vec3 struct {
	x, y, z int64
}

func (v vec3) add(v1 vec3) vec3 {
	return vec3{
		x: v.x + v1.x,
		y: v.y + v1.y,
		z: v.z + v1.z,
	}
}

func (v vec3) sub(v1 vec3) vec3 {
	return vec3{
		x: v.x - v1.x,
		y: v.y - v1.y,
		z: v.z - v1.z,
	}
}

func (v vec3) min(v1 vec3) vec3 {
	return vec3{
		x: min(v.x, v1.x),
		y: min(v.y, v1.y),
		z: min(v.z, v1.z),
	}
}

func (v vec3) max(v1 vec3) vec3 {
	return vec3{
		x: max(v.x, v1.x),
		y: max(v.y, v1.y),
		z: max(v.z, v1.z),
	}
}

func (v vec3) neighbors6() []vec3 {
	neighbors := make([]vec3, 6)
	for i, d := range []vec3{
		{1, 0, 0},
		{-1, 0, 0},
		{0, 1, 0},
		{0, -1, 0},
		{0, 0, 1},
		{0, 0, -1},
	} {
		neighbors[i] = v.add(d)
	}
	return neighbors
}

func (v vec3) neighbors26() []vec3 {
	neighbors := make([]vec3, 0, 26)
	for dx := int64(-1); dx <= 1; dx++ {
		for dy := int64(-1); dy <= 1; dy++ {
			for dz := int64(-1); dz <= 1; dz++ {
				if dx == 0 && dy == 0 && dz == 0 {
					continue
				}
				v1 := vec3{v.x + dx, v.y + dy, v.z + dz}
				neighbors = append(neighbors, v1)
			}
		}
	}
	return neighbors
}

func (v vec3) hamming(v1 vec3) int64 {
	return abs(v1.x-v.x) + abs(v1.y-v.y) + abs(v1.z-v.z)
}

type vec4 struct {
	x, y, z, w int64
}

func (v vec4) add(v1 vec4) vec4 {
	return vec4{
		x: v.x + v1.x,
		y: v.y + v1.y,
		z: v.z + v1.z,
		w: v.w + v1.w,
	}
}

func (v vec4) neighbors() []vec4 {
	neighbors := make([]vec4, 0, 80)
	for dx := int64(-1); dx <= 1; dx++ {
		for dy := int64(-1); dy <= 1; dy++ {
			for dz := int64(-1); dz <= 1; dz++ {
				for dw := int64(-1); dw <= 1; dw++ {
					if dx == 0 && dy == 0 && dz == 0 && dw == 0 {
						continue
					}
					v1 := vec4{v.x + dx, v.y + dy, v.z + dz, v.w + dw}
					neighbors = append(neighbors, v1)
				}
			}
		}
	}
	return neighbors
}

type grid[E any] struct {
	g    [][]E
	rows int64
	cols int64
}

func (g *grid[E]) init(rows, cols int64, v E) {
	if len(g.g) > 0 {
		panic("double init")
	}
	g.g = make([][]E, rows)
	for i := range g.g {
		row := make([]E, cols)
		for j := range row {
			row[j] = v
		}
		g.g[i] = row
	}
	g.rows = rows
	g.cols = cols
}

func (g *grid[E]) addRow(row []E) {
	if g.g == nil {
		g.cols = int64(len(row))
	} else if g.cols != int64(len(row)) {
		panic("non-rectangular grid")
	}
	g.g = append(g.g, row)
	g.rows++
}

func (g *grid[E]) contains(v vec2) bool {
	return v.x >= 0 && v.x < g.cols && v.y >= 0 && v.y < g.rows
}

func (g *grid[E]) at(v vec2) E {
	return g.g[v.y][v.x]
}

func (g *grid[E]) set(v vec2, e E) {
	g.g[v.y][v.x] = e
}

func (g *grid[E]) all() iter.Seq2[vec2, E] {
	return func(yield func(vec2, E) bool) {
		for y := int64(0); y < g.rows; y++ {
			for x := int64(0); x < g.cols; x++ {
				v := vec2{x, y}
				if !yield(v, g.at(v)) {
					return
				}
			}
		}
	}
}

func (g *grid[E]) vecs() iter.Seq[vec2] {
	return func(yield func(vec2) bool) {
		for y := int64(0); y < g.rows; y++ {
			for x := int64(0); x < g.cols; x++ {
				v := vec2{x, y}
				if !yield(v) {
					return
				}
			}
		}
	}
}

// ball gives all the grid positions that are part of the open ball with the
// given center and radius. The iterator gives the elements and their distance
// from the center.
func (g *grid[E]) ball(center vec2, radius int64) iter.Seq2[vec2, int64] {
	return func(yield func(vec2, int64) bool) {
		for d := int64(0); d < radius; d++ {
			v := center.add(vec2{0, -d})
			for _, dv := range []vec2{{-1, 1}, {1, 1}, {1, -1}, {-1, -1}} {
				for range d {
					if g.contains(v) {
						if !yield(v, d) {
							return
						}
					}
					v = v.add(dv)
				}
			}
		}
	}
}

func (g *grid[E]) insertCol(x int64, e E) {
	g.cols++
	for y := int64(0); y < g.rows; y++ {
		g.g[y] = slices.Insert(g.g[y], int(x), e)
	}
}

func (g *grid[E]) insertRow(y int64, e E) {
	g.rows++
	row := make([]E, g.cols)
	for y := range row {
		row[y] = e
	}
	g.g = slices.Insert(g.g, int(y), row)
}

func (g *grid[E]) clone() *grid[E] {
	g1 := &grid[E]{
		g:    make([][]E, len(g.g)),
		rows: g.rows,
		cols: g.cols,
	}
	for i, row := range g.g {
		g1.g[i] = slices.Clone(row)
	}
	return g1
}

func byteGridString(g *grid[byte]) string {
	var b strings.Builder
	for _, row := range g.g {
		b.Write(row)
		fmt.Fprintln(&b)
	}
	return b.String()
}

// Extra slice stuff not in slices.

func SliceMin[E constraints.Ordered](x []E) E {
	if len(x) == 0 {
		panic("SliceMin of 0 elements")
	}
	min := x[0]
	for i := 1; i < len(x); i++ {
		if x[i] < min {
			min = x[i]
		}
	}
	return min
}

func SliceMax[E constraints.Ordered](x []E) E {
	if len(x) == 0 {
		panic("SliceMax of 0 elements")
	}
	max := x[0]
	for i := 1; i < len(x); i++ {
		if x[i] > max {
			max = x[i]
		}
	}
	return max
}

func SliceSum[S ~[]E, E number](s S) E {
	var sum E
	for _, e := range s {
		sum += e
	}
	return sum
}

func SliceMap[S ~[]E1, E1, E2 any](s S, fn func(E1) E2) []E2 {
	r := make([]E2, len(s))
	for i, e1 := range s {
		r[i] = fn(e1)
	}
	return r
}

func SliceReduce[S ~[]E, E, R any](s S, initial R, fn func(R, E) R) R {
	r := initial
	for _, e := range s {
		r = fn(r, e)
	}
	return r
}

func SliceReverse[S ~[]E, E any](s S) {
	for i := 0; i < len(s)/2; i++ {
		j := len(s) - 1 - i
		s[i], s[j] = s[j], s[i]
	}
}

func StackPush[S ~[]E, E any](stk *S, e E) {
	*stk = append(*stk, e)
}

func StackPop[S ~[]E, E any](stk *S) E {
	e := (*stk)[len(*stk)-1]
	*stk = (*stk)[:len(*stk)-1]
	return e
}

func SlicePop[S ~[]E, E any](s *S) E {
	e := (*s)[0]
	*s = (*s)[1:]
	return e
}

func ParDo[S ~[]E, E any](s S, fn func(E)) {
	var wg sync.WaitGroup
	var idx atomic.Int64
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				i := int(idx.Add(1) - 1)
				if i >= len(s) {
					return
				}
				fn(s[i])
			}
		}()
	}
	wg.Wait()
}

var _ = pretty.Println // avoid requiring go mod tidy as I add/remove this dep
