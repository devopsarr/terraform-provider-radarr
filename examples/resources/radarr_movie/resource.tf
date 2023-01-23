resource "radarr_movie" "example" {
  monitored            = false
  title                = "The Matrix"
  path                 = "/movies/The_Matrix_1999"
  quality_profile_id   = 1
  tmdb_id              = 603
  minimum_availability = "inCinemas"
}