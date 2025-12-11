package main

import (
	"fmt"
	"math/rand"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	ScreenWidth  = 1200
	ScreenHeight = 800
	GridSize     = 40
	WallDensity  = 0.3
)

// Cell представляет клетку лабиринта
type Cell struct {
	X, Y    int
	Visited bool
	IsWall  bool
}

// LevelSize размер уровня
type LevelSize struct {
	Width, Height        int
	MinWidth, MaxWidth   int
	MinHeight, MaxHeight int
}

// Die представляет кубик с отслеживанием всех сторон
type Die struct {
	Top, Bottom, Front, Back, Left, Right int
	CurrentTop                            int
}

// Direction направление движения
type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

// Player представляет игрока-кубик
type Player struct {
	X, Y int
	Die  Die
}

// Level представляет уровень игры
type Level struct {
	Player Player
	Finish struct {
		X, Y   int
		Number int
	}
	Cells [][]Cell
	Size  LevelSize
	Won   bool
}

// NewDie создает новый кубик
func NewDie() Die {
	return Die{
		Top:        1,
		Bottom:     6,
		Front:      2,
		Back:       5,
		Left:       3,
		Right:      4,
		CurrentTop: 1,
	}
}

// Roll перекатывает кубик в указанном направлении
func (d *Die) Roll(dir Direction) {
	switch dir {
	case Up:
		newTop := d.Back
		d.Back = d.Bottom
		d.Bottom = d.Front
		d.Front = d.Top
		d.Top = newTop
	case Down:
		newTop := d.Front
		d.Front = d.Bottom
		d.Bottom = d.Back
		d.Back = d.Top
		d.Top = newTop
	case Left:
		newTop := d.Right
		d.Right = d.Bottom
		d.Bottom = d.Left
		d.Left = d.Top
		d.Top = newTop
	case Right:
		newTop := d.Left
		d.Left = d.Bottom
		d.Bottom = d.Right
		d.Right = d.Top
		d.Top = newTop
	}
	d.CurrentTop = d.Top
}

// GetDieColor возвращает цвет для числа на кубике
func GetDieColor(number int) rl.Color {
	switch number {
	case 1:
		return rl.Red
	case 2:
		return rl.Orange
	case 3:
		return rl.Yellow
	case 4:
		return rl.Green
	case 5:
		return rl.Blue
	case 6:
		return rl.Purple
	default:
		return rl.Gray
	}
}

// NewPlayer создает нового игрока
func NewPlayer(x, y int) Player {
	return Player{
		X:   x,
		Y:   y,
		Die: NewDie(),
	}
}

// Move двигает игрока в указанном направлении
func (p *Player) Move(dx, dy int, dir Direction) {
	p.X += dx
	p.Y += dy
	p.Die.Roll(dir)
}

// IsValidMove проверяет, можно ли двигаться в указанную клетку
func (l *Level) IsValidMove(x, y int) bool {
	// Проверяем границы сетки
	if x < 0 || x >= l.Size.Width || y < 0 || y >= l.Size.Height {
		return false
	}

	// Проверяем, не является ли клетка стеной
	if l.Cells[y][x].IsWall {
		return false
	}

	return true
}

// GenerateMaze генерирует лабиринт с помощью алгоритма Recursive Backtracking
// GenerateMaze генерирует лабиринт с гарантированным путем от старта к финишу
func (l *Level) GenerateMaze() {
	// Инициализация клеток
	l.Cells = make([][]Cell, l.Size.Height)
	for y := 0; y < l.Size.Height; y++ {
		l.Cells[y] = make([]Cell, l.Size.Width)
		for x := 0; x < l.Size.Width; x++ {
			l.Cells[y][x] = Cell{
				X:       x,
				Y:       y,
				Visited: false,
				IsWall:  false, // Начинаем со всеми проходимыми клетками
			}
		}
	}

	// Добавляем внешние стены
	for x := 0; x < l.Size.Width; x++ {
		l.Cells[0][x].IsWall = true
		l.Cells[l.Size.Height-1][x].IsWall = true
	}
	for y := 0; y < l.Size.Height; y++ {
		l.Cells[y][0].IsWall = true
		l.Cells[y][l.Size.Width-1].IsWall = true
	}

	// Добавляем внутренние стены (лабиринт)
	// Создаем несколько вертикальных и горизонтальных стен
	for y := 2; y < l.Size.Height-2; y += 3 {
		for x := 1; x < l.Size.Width-1; x++ {
			if rand.Float64() < WallDensity { // вероятность стены
				l.Cells[y][x].IsWall = true
			}
		}
	}

	for x := 2; x < l.Size.Width-2; x += 3 {
		for y := 1; y < l.Size.Height-1; y++ {
			if rand.Float64() < WallDensity { // вероятность стены
				l.Cells[y][x].IsWall = true
			}
		}
	}

	// Создаем гарантированный путь от старта к финишу
	l.createGuaranteedPath()

	// Добавляем случайные открытые проходы для соединения областей
	l.connectIsolatedAreas()

	// Создаем хотя бы одну 2x2 открытую область
	l.createOpenSpace2x2()

	// Убедимся, что старт и финиш проходимы
	l.Cells[0][0].IsWall = false
	l.Cells[l.Finish.Y][l.Finish.X].IsWall = false
}

// createGuaranteedPath создает гарантированный путь от старта к финишу
func (l *Level) createGuaranteedPath() {
	// Алгоритм для создания пути
	// Начинаем от старта (0,0) и идем к финишу
	x, y := 0, 0
	targetX, targetY := l.Finish.X, l.Finish.Y

	// Основное направление движения
	for x < targetX || y < targetY {
		// Решаем, двигаться ли вправо или вниз
		if x < targetX && (y >= targetY || rand.Float64() < 0.5) {
			// Двигаемся вправо
			for dx := 0; dx < 2 && x+dx < l.Size.Width; dx++ {
				l.Cells[y][x+dx].IsWall = false
			}
			x++
		} else if y < targetY {
			// Двигаемся вниз
			for dy := 0; dy < 2 && y+dy < l.Size.Height; dy++ {
				l.Cells[y+dy][x].IsWall = false
			}
			y++
		}
	}
}

// EnsureConnectivity проверяет и улучшает связность лабиринта
func (l *Level) EnsureConnectivity() {
	// Помечаем все клетки как непосещенные
	for y := 0; y < l.Size.Height; y++ {
		for x := 0; x < l.Size.Width; x++ {
			l.Cells[y][x].Visited = false
		}
	}

	// Проверяем доступность от старта
	queue := []struct{ x, y int }{{0, 0}}
	l.Cells[0][0].Visited = true

	// BFS для проверки связности
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		// Проверяем соседей
		directions := []struct{ dx, dy int }{{0, -1}, {0, 1}, {-1, 0}, {1, 0}}
		for _, dir := range directions {
			nx, ny := current.x+dir.dx, current.y+dir.dy
			if nx >= 0 && nx < l.Size.Width && ny >= 0 && ny < l.Size.Height {
				if !l.Cells[ny][nx].IsWall && !l.Cells[ny][nx].Visited {
					l.Cells[ny][nx].Visited = true
					queue = append(queue, struct{ x, y int }{nx, ny})
				}
			}
		}
	}

	// Если финиш не достижим, создаем путь
	if !l.Cells[l.Finish.Y][l.Finish.X].Visited {
		l.createDirectPathToFinish()
	}
}

// createDirectPathToFinish создает прямой путь к финишу
func (l *Level) createDirectPathToFinish() {
	// Создаем L-образный путь от старта к финишу
	// Сначала двигаемся по горизонтали, затем по вертикали
	for x := 0; x < l.Size.Width; x++ {
		l.Cells[0][x].IsWall = false
	}
	for y := 0; y < l.Size.Height; y++ {
		l.Cells[y][l.Size.Width-1].IsWall = false
	}
}

// connectIsolatedAreas соединяет изолированные области
func (l *Level) connectIsolatedAreas() {
	// Делаем дополнительные проходы в случайных местах
	for i := 0; i < l.Size.Width*l.Size.Height/10; i++ {
		x := rand.Intn(l.Size.Width-2) + 1
		y := rand.Intn(l.Size.Height-2) + 1

		// Делаем крестообразный проход
		for dy := -1; dy <= 1; dy++ {
			for dx := -1; dx <= 1; dx++ {
				if dx == 0 || dy == 0 { // Только вертикальные и горизонтальные
					nx, ny := x+dx, y+dy
					if nx >= 0 && nx < l.Size.Width && ny >= 0 && ny < l.Size.Height {
						if rand.Float64() < 0.5 {
							l.Cells[ny][nx].IsWall = false
						}
					}
				}
			}
		}
	}
}

// createOpenSpace2x2 создает как минимум одну открытую область 2x2
func (l *Level) createOpenSpace2x2() {
	// Выбираем случайную позицию для открытой области
	// Оставляем место для стен по краям
	x := rand.Intn(l.Size.Width-4) + 2
	y := rand.Intn(l.Size.Height-4) + 2

	// Создаем область 2x2 без стен
	for dy := 0; dy < 2; dy++ {
		for dx := 0; dx < 2; dx++ {
			if y+dy < l.Size.Height && x+dx < l.Size.Width {
				l.Cells[y+dy][x+dx].IsWall = false
			}
		}
	}

	// Обеспечиваем доступ к этой области, убирая стены вокруг нее
	for dy := -1; dy <= 2; dy++ {
		for dx := -1; dx <= 2; dx++ {
			nx, ny := x+dx, y+dy
			if nx >= 0 && nx < l.Size.Width && ny >= 0 && ny < l.Size.Height {
				// Убираем стены по периметру области 2x2
				if (dx == -1 || dx == 2 || dy == -1 || dy == 2) && rand.Float64() < 0.7 {
					l.Cells[ny][nx].IsWall = false
				}
			}
		}
	}
}

// NewLevel создает новый уровень с лабиринтом
func NewLevel(size LevelSize) Level {
	l := Level{}
	l.Size = size

	// Создаем игрока в левом верхнем углу
	l.Player = NewPlayer(0, 0)

	// Устанавливаем финиш в правом нижнем углу
	l.Finish.X = size.Width - 1
	l.Finish.Y = size.Height - 1
	l.Finish.Number = rand.Intn(6) + 1 // случайное число от 1 до 6

	// Генерируем лабиринт
	l.GenerateMaze()

	// Убедимся, что лабиринт связан
	l.EnsureConnectivity()

	// Убедимся, что старт и финиш проходимы
	l.Cells[0][0].IsWall = false
	l.Cells[l.Finish.Y][l.Finish.X].IsWall = false

	l.Won = false
	return l
}

// CheckWin проверяет условие победы
func (l *Level) CheckWin() bool {
	return l.Player.X == l.Finish.X &&
		l.Player.Y == l.Finish.Y &&
		l.Player.Die.CurrentTop == l.Finish.Number
}

// DrawMazeWalls рисует стены лабиринта как полные клетки
func DrawMazeWalls(cells [][]Cell, gridSize int, offsetX, offsetY int) {
	for y := 0; y < len(cells); y++ {
		for x := 0; x < len(cells[y]); x++ {
			if cells[y][x].IsWall {
				cellX := offsetX + x*gridSize
				cellY := offsetY + y*gridSize
				// Рисуем стену как закрашенную клетку
				rl.DrawRectangle(int32(cellX), int32(cellY), int32(gridSize), int32(gridSize), rl.DarkBrown)
				rl.DrawRectangleLines(int32(cellX), int32(cellY), int32(gridSize), int32(gridSize), rl.Black)
			}
		}
	}
}

// DrawDieWithSides рисует кубик с визуализацией всех сторон
func DrawDieWithSides(x, y int, die Die) {
	size := GridSize
	padding := 5
	dieSize := size - 2*padding

	drawX := x + padding
	drawY := y + padding

	// Рисуем основной кубик
	rl.DrawRectangle(int32(drawX), int32(drawY), int32(dieSize), int32(dieSize), GetDieColor(die.CurrentTop))
	rl.DrawRectangleLines(int32(drawX), int32(drawY), int32(dieSize), int32(dieSize), rl.Black)

	// Рисуем число на верхней стороне
	text := fmt.Sprintf("%d", die.CurrentTop)
	fontSize := int32(24)
	textWidth := rl.MeasureText(text, fontSize)
	textX := drawX + (dieSize-int(textWidth))/2
	textY := drawY + (dieSize-24)/2
	rl.DrawText(text, int32(textX), int32(textY), fontSize, rl.White)

	// Рисуем стороны кубика как цветные полоски
	stripHeight := 4
	margin := 2

	// Левая полоска (Left сторона)
	rl.DrawRectangle(int32(drawX-margin-stripHeight), int32(drawY+margin), int32(stripHeight), int32(dieSize-2*margin), GetDieColor(die.Left))

	// Правая полоска (Right сторона)
	rl.DrawRectangle(int32(drawX+dieSize+margin), int32(drawY+margin), int32(stripHeight), int32(dieSize-2*margin), GetDieColor(die.Right))

	// Передняя полоска (Front сторона) - внизу
	rl.DrawRectangle(int32(drawX+margin), int32(drawY+dieSize+margin+2), int32(dieSize-2*margin), int32(stripHeight), GetDieColor(die.Front))

	// Задняя полоска (Back сторона) - вверху
	rl.DrawRectangle(int32(drawX+margin), int32(drawY-margin-stripHeight-2), int32(dieSize-2*margin), int32(stripHeight), GetDieColor(die.Back))

	// Подписи к полоскам
	fontSizeSmall := int32(12)
	rl.DrawText(fmt.Sprintf("%d", die.Left), int32(drawX-margin-stripHeight-15), int32(drawY+dieSize/2-6), fontSizeSmall, rl.Black)
	rl.DrawText(fmt.Sprintf("%d", die.Right), int32(drawX+dieSize+margin+stripHeight+2), int32(drawY+dieSize/2-6), fontSizeSmall, rl.Black)
	rl.DrawText(fmt.Sprintf("%d", die.Front), int32(drawX+dieSize/2-6), int32(drawY+dieSize+margin+stripHeight+5), fontSizeSmall, rl.Black)
	rl.DrawText(fmt.Sprintf("%d", die.Back), int32(drawX+dieSize/2-6), int32(drawY-margin-stripHeight-2-stripHeight-5), fontSizeSmall, rl.Black)
}

// DrawGrid рисует фон сетки
func DrawGrid(width, height, gridSize int, offsetX, offsetY int) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			var color rl.Color
			if (x+y)%2 == 0 {
				color = rl.LightGray
			} else {
				color = rl.Gray
			}

			cellX := offsetX + x*gridSize
			cellY := offsetY + y*gridSize

			rl.DrawRectangle(int32(cellX), int32(cellY), int32(gridSize), int32(gridSize), color)
			rl.DrawRectangleLines(int32(cellX), int32(cellY), int32(gridSize), int32(gridSize), rl.DarkGray)
		}
	}
}

// DrawFinish рисует финишную клетку
func DrawFinish(x, y, gridSize, offsetX, offsetY int, number int) {
	cellX := offsetX + x*gridSize
	cellY := offsetY + y*gridSize

	// Рисуем финишную клетку
	rl.DrawRectangle(int32(cellX), int32(cellY), int32(gridSize), int32(gridSize), rl.Gold)
	rl.DrawRectangleLines(int32(cellX), int32(cellY), int32(gridSize), int32(gridSize), rl.Black)

	// Рисуем число на финише
	text := fmt.Sprintf("%d", number)
	fontSize := int32(24)
	textWidth := rl.MeasureText(text, fontSize)
	textX := cellX + (gridSize-int(textWidth))/2
	textY := cellY + (gridSize-24)/2

	rl.DrawText(text, int32(textX), int32(textY), fontSize, rl.Black)
}

// DrawLevelSizeUI рисует UI для выбора размера уровня
func DrawLevelSizeUI(selectedWidth, selectedHeight int) {
	// Фон для UI
	rl.DrawRectangle(10, 10, 300, 120, rl.White)
	rl.DrawRectangleLines(10, 10, 300, 120, rl.Black)

	// Заголовок
	rl.DrawText("Level Size:", 20, 20, 20, rl.Black)

	// Размеры
	widthText := fmt.Sprintf("Width: %d", selectedWidth)
	heightText := fmt.Sprintf("Height: %d", selectedHeight)

	rl.DrawText(widthText, 20, 50, 18, rl.Black)
	rl.DrawText(heightText, 20, 75, 18, rl.Black)

	// Инструкции
	rl.DrawText("1-9: Width  |  Q-I: Height", 20, 100, 14, rl.DarkGray)
	rl.DrawText("R: Regenerate  |  Enter: Start", 20, 115, 14, rl.DarkGray)
}

// DrawUI рисует пользовательский интерфейс
func DrawUI(level Level, gridSize, offsetX, offsetY int) {
	// Информация о текущем числе кубика
	currentText := fmt.Sprintf("Current: %d", level.Player.Die.CurrentTop)
	rl.DrawText(currentText, 320, 20, 24, rl.Black)

	// Информация о числе на финише
	targetText := fmt.Sprintf("Target: %d", level.Finish.Number)
	rl.DrawText(targetText, 320, 50, 24, rl.Black)

	// Позиция игрока
	posText := fmt.Sprintf("Pos: (%d,%d)", level.Player.X, level.Player.Y)
	rl.DrawText(posText, 320, 80, 18, rl.DarkGray)

	// Размер уровня
	sizeText := fmt.Sprintf("Size: %dx%d", level.Size.Width, level.Size.Height)
	rl.DrawText(sizeText, 320, 105, 18, rl.DarkGray)

	// Инструкции
	instructions := "WASD/Arrows: Move | R: Regenerate | 1-9/Q-I: Size"
	rl.DrawText(instructions, 320, 130, 16, rl.DarkGray)

	// Сообщение о победе
	if level.Won {
		winText := "YOU WIN! Press R for new level"
		textWidth := rl.MeasureText(winText, 30)
		textX := (ScreenWidth - int(textWidth)) / 2
		textY := ScreenHeight/2 - 15

		rl.DrawRectangle(int32(textX-10), int32(textY-10), int32(textWidth+20), 60, rl.Green)
		rl.DrawRectangleLines(int32(textX-10), int32(textY-10), int32(textWidth+20), 60, rl.Black)
		rl.DrawText(winText, int32(textX), int32(textY), 30, rl.White)
	}
}

// HandleInput обрабатывает ввод игрока
func HandleInput(level *Level, gridSize, offsetX, offsetY *int, currentSize *LevelSize) {
	// Изменение размера уровня
	if rl.IsKeyPressed(rl.KeyOne) {
		currentSize.Width = 10
	}
	if rl.IsKeyPressed(rl.KeyTwo) {
		currentSize.Width = 12
	}
	if rl.IsKeyPressed(rl.KeyThree) {
		currentSize.Width = 15
	}
	if rl.IsKeyPressed(rl.KeyFour) {
		currentSize.Width = 18
	}
	if rl.IsKeyPressed(rl.KeyFive) {
		currentSize.Width = 20
	}
	if rl.IsKeyPressed(rl.KeySix) {
		currentSize.Width = 25
	}
	if rl.IsKeyPressed(rl.KeySeven) {
		currentSize.Width = 30
	}
	if rl.IsKeyPressed(rl.KeyEight) {
		currentSize.Width = 35
	}
	if rl.IsKeyPressed(rl.KeyNine) {
		currentSize.Width = 40
	}

	if rl.IsKeyPressed(rl.KeyQ) {
		currentSize.Height = 8
	}
	if rl.IsKeyPressed(rl.KeyW) {
		currentSize.Height = 10
	}
	if rl.IsKeyPressed(rl.KeyE) {
		currentSize.Height = 12
	}
	if rl.IsKeyPressed(rl.KeyT) {
		currentSize.Height = 15
	}
	if rl.IsKeyPressed(rl.KeyY) {
		currentSize.Height = 18
	}
	if rl.IsKeyPressed(rl.KeyU) {
		currentSize.Height = 20
	}
	if rl.IsKeyPressed(rl.KeyI) {
		currentSize.Height = 25
	}

	// Перегенерация уровня
	if rl.IsKeyPressed(rl.KeyR) {
		*level = NewLevel(*currentSize)
		*gridSize = GridSize
		*offsetX = (ScreenWidth - currentSize.Width*GridSize) / 2
		*offsetY = (ScreenHeight - currentSize.Height*GridSize) / 2
		return
	}

	if level.Won {
		return
	}

	// Движение по WASD
	if rl.IsKeyPressed(rl.KeyW) || rl.IsKeyPressed(rl.KeyUp) {
		newX, newY := level.Player.X, level.Player.Y-1
		if level.IsValidMove(newX, newY) {
			level.Player.Move(0, -1, Up)
		}
	}
	if rl.IsKeyPressed(rl.KeyS) || rl.IsKeyPressed(rl.KeyDown) {
		newX, newY := level.Player.X, level.Player.Y+1
		if level.IsValidMove(newX, newY) {
			level.Player.Move(0, 1, Down)
		}
	}
	if rl.IsKeyPressed(rl.KeyA) || rl.IsKeyPressed(rl.KeyLeft) {
		newX, newY := level.Player.X-1, level.Player.Y
		if level.IsValidMove(newX, newY) {
			level.Player.Move(-1, 0, Left)
		}
	}
	if rl.IsKeyPressed(rl.KeyD) || rl.IsKeyPressed(rl.KeyRight) {
		newX, newY := level.Player.X+1, level.Player.Y
		if level.IsValidMove(newX, newY) {
			level.Player.Move(1, 0, Right)
		}
	}

	// Проверяем победу после движения
	if level.CheckWin() {
		level.Won = true
	}
}

func main() {
	// Инициализируем генератор случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Создаем окно
	rl.InitWindow(ScreenWidth, ScreenHeight, "KubeGame - Labyrinth Die Puzzle")
	rl.SetTargetFPS(60)

	// Настройки уровня
	currentSize := LevelSize{
		Width:     15,
		Height:    10,
		MinWidth:  5,
		MaxWidth:  50,
		MinHeight: 5,
		MaxHeight: 40,
	}

	// Создаем уровень
	level := NewLevel(currentSize)

	// Вычисляем позиционирование
	gridSize := GridSize
	offsetX := (ScreenWidth - currentSize.Width*gridSize) / 2
	offsetY := (ScreenHeight - currentSize.Height*gridSize) / 2

	// Главный игровой цикл
	for !rl.WindowShouldClose() {
		// Обновление
		HandleInput(&level, &gridSize, &offsetX, &offsetY, &currentSize)

		// Рендеринг
		rl.BeginDrawing()
		rl.ClearBackground(rl.RayWhite)

		// Рисуем игровое поле
		DrawGrid(level.Size.Width, level.Size.Height, gridSize, offsetX, offsetY)
		DrawMazeWalls(level.Cells, gridSize, offsetX, offsetY)
		DrawFinish(level.Finish.X, level.Finish.Y, gridSize, offsetX, offsetY, level.Finish.Number)

		// Рисуем игрока
		playerX := offsetX + level.Player.X*gridSize
		playerY := offsetY + level.Player.Y*gridSize
		DrawDieWithSides(playerX, playerY, level.Player.Die)

		// Рисуем UI
		DrawLevelSizeUI(currentSize.Width, currentSize.Height)
		DrawUI(level, gridSize, offsetX, offsetY)

		rl.EndDrawing()
	}

	// Закрываем окно
	rl.CloseWindow()
}
