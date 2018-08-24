library(tidyverse)

conv_files <- Sys.glob(paste("*", "oldest_states_hist", "csv", sep = "."))

## import state-value hist
conv <- c()
for (conv_file in conv_files) {
  name <- unlist(strsplit(conv_file, split = "\\."))[1]
  conv <-
    rbind(
      conv,
      read.csv(conv_file, stringsAsFactors = FALSE, header = FALSE) %>%
        select(state = V1, time = V2, value = V3) %>%
        mutate(name = name)
    )
}
conv$state <- as.character(conv$state)

plot <-
  ggplot(conv) +
  geom_line(aes(x = time, y = value, color = state), size = 1) +
  facet_grid(name ~ .)
print(plot)
