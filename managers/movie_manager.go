package managers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"movie_management/models"
	"movie_management/producer"
	"movie_management/request"
	"movie_management/response"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

// type MovieProcessingConfig struct {
// 	EnableProducer bool
// }

// var movieConfig = MovieProcessingConfig{EnableProducer: false}

func CreateMovie(req request.MovieRequest) (response.MovieResponse, error) {
	o := orm.NewOrm()
	log.Println("reqmanagers----->")
	existingMovie := models.Movie{}
	err := o.QueryTable(&models.Movie{}).Filter("title", req.Title).One(&existingMovie)
	if err == nil {
		return response.MovieResponse{}, fmt.Errorf("movie with title '%s' already exists", req.Title)
	} else if err != orm.ErrNoRows {
		return response.MovieResponse{}, err
	}
	var movie = models.Movie{
		Title:     req.Title,
		Genre:     req.Genre,
		Year:      req.Year,
		Rating:    req.Rating,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	log.Println("reqmanagers 2----->")

	// _, err = o.Insert(&movie)
	// if err != nil {
	// 	return response.MovieResponse{}, err
	// }

	jsonData, err := json.Marshal(movie)
	if err != nil {
		return response.MovieResponse{}, err
	}

	rmp := producer.NewProducer()
	producerService := producer.NewProducerService(rmp)
	if err := producerService.Initialize(); err != nil {
		log.Println("Failed to initialize producer service:", err)
		return response.MovieResponse{}, err
	}

	if err := producerService.Publish(jsonData); err != nil {
		log.Println("Error publishing movie data:", err)
		return response.MovieResponse{}, err
	}

	return response.MovieResponse{
		ID:        movie.Id,
		Title:     movie.Title,
		Genre:     movie.Genre,
		Year:      movie.Year,
		Rating:    movie.Rating,
		CreatedAt: movie.CreatedAt,
		UpdatedAt: movie.UpdatedAt,
	}, nil
}

func UpdateMovie(id int, req request.MovieRequest) (response.MovieResponse, error) {
	o := orm.NewOrm()

	movie := models.Movie{Id: id}
	err := o.Read(&movie)
	if err != nil {
		return response.MovieResponse{}, fmt.Errorf("movie not found")
	}

	movie.Title = req.Title
	movie.Genre = req.Genre
	movie.Year = req.Year
	movie.Rating = req.Rating
	movie.UpdatedAt = time.Now()

	_, err = o.Update(&movie)
	if err != nil {
		return response.MovieResponse{}, fmt.Errorf("failed to update movie: %v", err)
	}
	responseMovie := &response.MovieResponse{
		ID:        movie.Id,
		Title:     movie.Title,
		Genre:     movie.Genre,
		Year:      movie.Year,
		Rating:    movie.Rating,
		CreatedAt: movie.CreatedAt,
		UpdatedAt: movie.UpdatedAt,
	}

	return *responseMovie, nil
}

func DeleteMovie(id int) error {
	o := orm.NewOrm()
	existingMovie := models.Movie{Id: id}
	if err := o.Read(&existingMovie); err != nil {
		if err == orm.ErrNoRows {
			return errors.New("movie not found")
		}
		return err
	}
	if _, err := o.Delete(&existingMovie); err != nil {
		return err
	}
	return nil
}

func ListMovies(o orm.Ormer, req request.Req) (response.ListMoviesResponse, error) {
	var movies []response.MovieResponse

	query := o.QueryTable(&models.Movie{})

	if req.Filter != "" {
		query = query.Filter("title__icontains", req.Filter).Filter("genre__icontains", req.Filter)
	}

	if req.Year != 0 {
		query = query.Filter("year", req.Year)
	}

	if req.OrderBy != "" && req.Order != "" {
		order := map[string]string{"asc": "", "desc": "-"}
		if orderDirection, ok := order[req.Order]; ok {
			query = query.OrderBy(orderDirection + req.OrderBy)
		} else {
			return response.ListMoviesResponse{}, fmt.Errorf("invalid order direction: %s", req.Order)
		}
	}

	offset := (req.PageNo - 1) * req.PageSize
	query = query.Limit(req.PageSize).Offset(offset)

	if _, err := query.All(&movies); err != nil {
		return response.ListMoviesResponse{}, fmt.Errorf("failed to fetch movies: %v", err)
	}

	countQuery := o.QueryTable(&models.Movie{})
	if req.Filter != "" {
		countQuery = countQuery.Filter("title__icontains", req.Filter).Filter("genre__icontains", req.Filter)
	}
	if req.Year != 0 {
		countQuery = countQuery.Filter("year", req.Year)
	}

	total, err := countQuery.Count()
	if err != nil {
		return response.ListMoviesResponse{}, fmt.Errorf("failed to count movies: %v", err)
	}

	lastPage := (total + int64(req.PageSize) - 1) / int64(req.PageSize)
	if lastPage == 0 {
		lastPage = 1
	}

	return response.ListMoviesResponse{
		Movies:      movies,
		PageNo:      req.PageNo,
		PageSize:    req.PageSize,
		TotalCount:  int(total),
		LastPage:    int(lastPage),
		CurrentPage: req.PageNo,
	}, nil
}

func GetMoviesById(id int) (response.MovieResponse, error) {
	o := orm.NewOrm()

	var movie models.Movie

	if err := o.QueryTable(new(models.Movie)).Filter("Id", id).One(&movie); err != nil {
		if err == orm.ErrNoRows {
			return response.MovieResponse{}, fmt.Errorf("movie not found")
		}
		return response.MovieResponse{}, fmt.Errorf("failed to fetch movie: %v", err)
	}

	movieResponse := response.MovieResponse{
		ID:        movie.Id,
		Title:     movie.Title,
		Genre:     movie.Genre,
		Year:      movie.Year,
		Rating:    movie.Rating,
		CreatedAt: movie.CreatedAt,
		UpdatedAt: movie.UpdatedAt,
	}

	return movieResponse, nil
}

func GetMovieAnalytics() (map[string]interface{}, error) {
	genreCounts, err := fetchGenreCounts()
	if err != nil {
		return nil, err
	}

	topRatedMovies, err := fetchTopRatedMovies()
	if err != nil {
		return nil, err
	}

	recentlyAddedMovies, err := fetchRecentlyAddedMovies()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"genreCounts":         genreCounts,
		"topRatedMovies":      topRatedMovies,
		"recentlyAddedMovies": recentlyAddedMovies,
	}, nil
}

func fetchGenreCounts() ([]response.GenreCount, error) {
	o := orm.NewOrm()
	var genreCounts []response.GenreCount

	if _, err := o.Raw("SELECT genre, COUNT(*) AS count FROM movie GROUP BY genre").QueryRows(&genreCounts); err != nil {
		return nil, fmt.Errorf("failed to fetch genre counts: %v", err)
	}
	return genreCounts, nil
}

func fetchTopRatedMovies() ([]response.MovieResponse, error) {
	o := orm.NewOrm()

	var maxRating float64
	if err := o.Raw("SELECT MAX(rating) FROM movie").QueryRow(&maxRating); err != nil {
		return nil, fmt.Errorf("failed to fetch max rating: %v", err)
	}

	var movies []response.MovieResponse
	if _, err := o.Raw("SELECT * FROM movie WHERE rating = ?", maxRating).QueryRows(&movies); err != nil {
		return nil, fmt.Errorf("failed to fetch top-rated movies: %v", err)
	}
	return movies, nil
}

func fetchRecentlyAddedMovies() ([]response.MovieResponse, error) {
	o := orm.NewOrm()
	var movies []response.MovieResponse

	timeOneMinuteAgo := time.Now().Add(-time.Minute)
	log.Printf("Fetching movies added since: %v", timeOneMinuteAgo)

	sql := "SELECT * FROM movie WHERE created_at >= ?"
	if _, err := o.Raw(sql, timeOneMinuteAgo).QueryRows(&movies); err != nil {
		return nil, fmt.Errorf("failed to fetch recently added movies: %v", err)
	}

	if len(movies) == 0 {
		log.Println("No movies found added in the last minute")
		return []response.MovieResponse{}, nil
	}

	return movies, nil
}
