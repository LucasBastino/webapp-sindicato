package page

type Pagination struct {
	SearchKey string
	Offset    int

	CurrentPage  int
	FirstPage    int
	PreviousPage int
	SomeBefore   int
	SixBefore    int
	FiveBefore   int
	FourBefore   int
	ThreeBefore  int
	TwoBefore    int
	TwoAfter     int
	ThreeAfter   int
	FourAfter    int
	FiveAfter    int
	SixAfter     int
	SomeAfter    int
	NextPage     int
	LastPage     int

	TotalPages      int
	TotalPagesArray []int
}

type PaginationMeta struct {
	TotalPages int
	Offset     int
	SomeBefore int
	SomeAfter  int
}

func calculatePaginationMeta(currentPage int, totalRows int) PaginationMeta {

	// setting totalPages
	var totalPages int
	// si no hay filas no hay paginas, se pone 1 para que calcule bien el offset
	if totalRows == 0 {
		totalPages = 1
		// si la cantidad de filas es un multiplo de 10 entran justo y no sobran
	} else if totalRows%15 == 0 {
		totalPages = totalRows / 15
		// sino no entran justo y se agrega una pagina mas
	} else {
		totalPages = (totalRows / 15) + 1
	}

	// setting currentPage and offset
	var offset int

	// PODES HACER UNA FUNCION DE ESTO O METER UN SWITCH
	//  si currentPage es menor a 1, currentPage ahora es 1 y muestra los primeros 20
	if currentPage <= 1 {
		offset = 0
	}

	// si currentPage es mayor a totalPages, currentPage ahora es totalPages
	// y muestra los ultimos members
	if currentPage > totalPages {
		currentPage = totalPages
		offset = (currentPage - 1) * 15
	}

	// si currentPage es mayor a 1, muestra los miembros calculando el offset * 15
	if currentPage > 1 {
		offset = (currentPage - 1) * 15
	}

	// setting aproximador
	someBefore := totalPages / 6
	someAfter := totalPages / 6
	// si se pasa de la ultima que te lleve a la ultima
	if someAfter+currentPage > totalPages {
		someAfter = totalPages - currentPage
		// si se pasa de la primera que te lleve a la primera
	} else if currentPage-someBefore < 1 {
		someBefore = currentPage - 1
	}
	return PaginationMeta{
		TotalPages: totalPages,
		Offset:     offset,
		SomeBefore: someBefore,
		SomeAfter:  someAfter,
	}
}

func getTotalPagesArray(totalPages int) []int {
	// devuelve el array para que se pueda recorrer en el template
	var totalPagesArray []int
	if totalPages <= 10 {
		for i := 1; i <= totalPages; i++ {
			totalPagesArray = append(totalPagesArray, i)
		}
	}
	return totalPagesArray
}

func setPagination(currentPage int, paginationMeta PaginationMeta, totalPagesArray []int) Pagination {
	return Pagination{
		CurrentPage:     currentPage,
		Offset:          paginationMeta.Offset,
		FirstPage:       1,
		PreviousPage:    currentPage - 1,
		SomeBefore:      currentPage - paginationMeta.SomeBefore,
		SixBefore:       currentPage - 6,
		FiveBefore:      currentPage - 5,
		FourBefore:      currentPage - 4,
		ThreeBefore:     currentPage - 3,
		TwoBefore:       currentPage - 2,
		TwoAfter:        currentPage + 2,
		ThreeAfter:      currentPage + 3,
		FourAfter:       currentPage + 4,
		FiveAfter:       currentPage + 5,
		SixAfter:        currentPage + 6,
		SomeAfter:       currentPage + paginationMeta.SomeAfter,
		NextPage:        currentPage + 1,
		LastPage:        paginationMeta.TotalPages,
		TotalPages:      paginationMeta.TotalPages,
		TotalPagesArray: totalPagesArray,
	}
}

func BuildPagination(currentPage, totalRows int) Pagination {
	// calcular totalPages
	paginationMeta := calculatePaginationMeta(currentPage, totalRows)

	// hago un array para poder recorrerlo y crear botones cuando hay menos de 10 paginas en el template
	totalPagesArray := getTotalPagesArray(paginationMeta.TotalPages)

	return setPagination(currentPage, paginationMeta, totalPagesArray)
}