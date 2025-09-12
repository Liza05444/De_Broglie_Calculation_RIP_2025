package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type Particle struct {
	ID          int
	Name        string
	Mass        float64
	Description string
	Image       string
}

var particles = []Particle{
	{
		ID:   1,
		Name: "Электрон",
		Mass: 9.109384e-31,
		Description: `В атоме электрон движется вокруг ядра, образуя электронные оболочки со скоростью, близкой к скорости света.
					 Количество электронов в оболочке определяет ее энергию и форму, а также определяет свойства атома, такие как химические свойства и атомный номер.
					 Электрон является фермионом, то есть частицей, подчиняющейся статистике Ферми-Дирака и имеющей полуцелое значение спина.
					 В физике используется как пример элементарной частицы, потому что он имеет наименьшую массу из всех известных частиц, и поэтому он может быть использован для изучения фундаментальных законов физики.
					 Изучение этой частицы помогает нам лучше понять, как работают атомы и молекулы, и как они взаимодействуют друг с другом.`,
		Image: "http://127.0.0.1:9000/particles/1.png",
	},
	{
		ID:   2,
		Name: "Протон",
		Mass: 1.672622e-27,
		Description: `Протон — одна из трёх (вместе с нейтроном и электроном) элементарных частиц, из которых построено обычное вещество.
					 Протоны входят в состав атомных ядер; порядковый номер химического элемента в таблице Менделеева равен количеству протонов в его ядре.
					 Спин протона равен 1/2, поэтому протоны подчиняются статистике Ферми-Дирака. 
					 У протона существует античастица – антипротон.
					 Название «протон» предложено Э. Резерфордом в 1920 году.`,
		Image: "http://127.0.0.1:9000/particles/2.png",
	},
	{
		ID:   3,
		Name: "Нейтрон",
		Mass: 1.674928e-27,
		Description: `Нейтрон (от лат. neuter — ни тот, ни другой) — тяжёлая субатомная частица, не имеющая электрического заряда.
					 Открыт в 1932 г. Дж. Чедвиком. Спин нейтрона равен 1/2.
					 В свободном состоянии нейтрон нестабилен – распадается на протон, электрон и антинейтрино;
					 время жизни составляет примерно 886 с.
					 Так как нейтрон электрически нейтрален, он легко проникает в атомные ядра при любой энергии и с большой вероятностью вызывает ядерные реакции.`,
		Image: "http://127.0.0.1:9000/particles/3.png",
	},
	{
		ID:   4,
		Name: "Альфа-частица",
		Mass: 6.644657e-27,
		Description: `Альфа-частица — положительно заряженная частица, образованная двумя протонами и двумя нейтронами; ядро атома гелия-4.
					 Впервые обнаружены Эрнестом Резерфордом в 1899 году и он же дал название этому виду излучения.
					 Альфа-частицы могут вызывать ядерные реакции; в первой искусственно вызванной ядерной реакции, проведённой Э. Резерфордом в 1919 году 
					 (превращение ядер азота в ядра кислорода) участвовали именно альфа-частицы.
					 Поток альфа-частиц называют альфа-лучами или альфа-излучением.`,
		Image: "http://127.0.0.1:9000/particles/4.png",
	},
	{
		ID:   5,
		Name: "Мюон",
		Mass: 1.883532e-28,
		Description: `Мюон — нестабильная заряженная элементарная частица со спином 1/2 и временем жизни 2,2 мкс. Мюон имеет античастицу.
					 Мюоны были открыты Карлом Андерсоном и Сетом Неддермайером в 1937 году, во время исследования космического излучения.
					 Они обнаружили частицы, которые при прохождении через магнитное поле отклонялись в меньшей степени, чем электроны, но сильнее, чем протоны.
					 Было сделано предположение, что они имеют элементарный заряд, и для объяснения различия в отклонении было необходимо,
					 чтобы эти частицы имели промежуточную массу, которая лежала бы между массами электрона и протона.`,
		Image: "http://127.0.0.1:9000/particles/5.png",
	},
	{
		ID:   6,
		Name: "Тау-лептон",
		Mass: 3.167541e-27,
		Description: `Тау-лептон, таон (от греческой буквы греч. τ — тау, использующейся для обозначения) — нестабильная элементарная частица
				     с отрицательным электрическим зарядом и спином 1/2. В Стандартной Модели физики элементарных частиц классифицируется как часть лептонного семейства фермионов (вместе с электроном, мюоном и нейтрино).
					 Как и все фундаментальные частицы, тау-лептон имеет античастицу с зарядом противоположного знака, но с равной массой и спином: антитау-лептон (антитаон).
					 Тау-лептон был открыт в 1975 году на электрон-позитронном коллайдере SPEAR в Национальной ускорительной лаборатории SLAC (Стэнфорд, США) М. Перлом и сотрудниками.
					 За открытие этой частицы Мартин Перл получил Нобелевскую премию по физике за 1995 год.`,
		Image: "http://127.0.0.1:9000/particles/6.png",
	},
	{
		ID:   7,
		Name: "Бозон Хиггса",
		Mass: 2.232122e-25,
		Description: `Бозон Хиггса — в современной теории элементарных частиц это неделимая частица, которая отвечает за механизм появления масс у некоторых других элементарных частиц.
					 В 1964 году британский физик Питер Хиггс вместе с другими учеными предположил, что существует особое поле, при взаимодействии с которым частицы приобретают массу.
					 Позже его назвали полем Хиггса, а процесс обретения массы — хиггсовским механизмом. Изучить, как работает этот процесс, можно только через измерения свойств хиггсовского бозона.
					 Без обнаружения бозона изучить это поле не удавалось.
					 Поэтому открытие бозона и понимание его свойств представлялось ученым важнейшей задачей.`,
		Image: "http://127.0.0.1:9000/particles/7.png",
	},
}

func (r *Repository) GetParticles() ([]Particle, error) {
	if len(particles) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return particles, nil
}

func (r *Repository) GetParticle(id int) (Particle, error) {
	particles, err := r.GetParticles()
	if err != nil {
		return Particle{}, err
	}

	for _, particle := range particles {
		if particle.ID == id {
			return particle, nil
		}
	}

	return Particle{}, fmt.Errorf("частица не найдена")
}

func (r *Repository) GetParticleByName(name string) ([]Particle, error) {
	particles, err := r.GetParticles()
	if err != nil {
		return []Particle{}, err
	}

	var result []Particle
	for _, particle := range particles {
		if strings.Contains(strings.ToLower(particle.Name), strings.ToLower(name)) {
			result = append(result, particle)
		}
	}

	return result, nil
}

type RequestDeBroglieCalculation struct {
	ID   int
	Name string
}

type DeBroglieCalculation struct {
	IDRequestDeBroglieCalculation int
	IDParticle                    int
	Speed                         float64
	Wavelength                    float64
}

var requestDeBroglieCalculations = map[RequestDeBroglieCalculation][]DeBroglieCalculation{
	{ID: 1, Name: "Эксперимент № 1"}: {
		{
			IDRequestDeBroglieCalculation: 1,
			IDParticle:                    particles[0].ID,
			Speed:                         1000000,
			Wavelength:                    7.274e-10,
		},
		{
			IDRequestDeBroglieCalculation: 1,
			IDParticle:                    particles[1].ID,
			Speed:                         100000,
			Wavelength:                    3.960e-12,
		},
		{
			IDRequestDeBroglieCalculation: 1,
			IDParticle:                    particles[2].ID,
			Speed:                         1000,
			Wavelength:                    3.956e-10,
		},
	},
	{ID: 2, Name: "Эксперимент № 2"}: {
		{
			IDRequestDeBroglieCalculation: 2,
			IDParticle:                    particles[2].ID,
			Speed:                         1000,
			Wavelength:                    3.956e-10,
		},
		{
			IDRequestDeBroglieCalculation: 2,
			IDParticle:                    particles[1].ID,
			Speed:                         100000,
			Wavelength:                    3.960e-12,
		},
	},
	{ID: 3, Name: "Эксперимент № 3"}: {
		{
			IDRequestDeBroglieCalculation: 3,
			IDParticle:                    particles[2].ID,
			Speed:                         1000,
			Wavelength:                    3.956e-10,
		},
	},
}

func (r *Repository) GetDeBroglieCalculationsForRequest(requestID int) (RequestDeBroglieCalculation, []DeBroglieCalculation, error) {
	for req, calculations := range requestDeBroglieCalculations {
		if req.ID == requestID {
			return req, calculations, nil
		}
	}
	return RequestDeBroglieCalculation{}, nil, fmt.Errorf("такой заявки нет")
}
