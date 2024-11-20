use actix_web::{post, web, HttpResponse, Responder};
use once_cell::sync::Lazy;
use regex::Regex;
use serde::{Deserialize, Serialize};
use std::collections::{HashMap, HashSet};

/// Request data for the similarity endpoint.
#[derive(Deserialize)]
struct SimilarityRequest {
    text1: String,
    text2: String,
}

/// Response data for the similarity endpoint.
#[derive(Serialize)]
struct SimilarityResponse {
    similarity: f64,
    interpretation: String,
}

/// Calculate the similarity between two texts.
#[post("/similarity")]
pub async fn similarity(data: web::Json<SimilarityRequest>) -> impl Responder {
    let normalized1 = normalize_text(&data.text1);
    let normalized2 = normalize_text(&data.text2);

    let words1: Vec<&str> = normalized1.split_whitespace().collect();
    let words2: Vec<&str> = normalized2.split_whitespace().collect();

    // Generate frequency maps for both texts
    let freq_map1 = generate_frequency_map(&words1);
    let freq_map2 = generate_frequency_map(&words2);

    // Create a vector of unique words
    let uniq: Vec<&str> = freq_map1.keys()
        .chain(freq_map2.keys())
        .cloned()
        .collect::<HashSet<&str>>()
        .into_iter()
        .collect();

    // Calculate term frequency for both texts
    let total1 = words1.len();
    let total2 = words2.len();

    let tf1 = calculate_tf(&uniq, &freq_map1, total1);
    let tf2 = calculate_tf(&uniq, &freq_map2, total2);

    // Calculate inverse document frequency
    let idf = calculate_idf(&uniq, &freq_map1, &freq_map2);

    // Calculate tf-idf for both texts
    let tf_idf1 = calculate_tf_idf(&tf1, &idf);
    let tf_idf2 = calculate_tf_idf(&tf2, &idf);

    // Calculate cosine similarity
    let similarity = calculate_similarity(&tf_idf1, &tf_idf2);

    // Round similarity to 3 decimal places
    let similarity = (similarity * 1000.0).round() / 1000.0;
    let interpretation = interpret_similarity(similarity);

    // Return the similarity as JSON
    HttpResponse::Ok().json(SimilarityResponse { similarity, interpretation })
}

/// Normalize text by converting to lowercase, removing punctuation, and collapsing whitespace.
fn normalize_text(text: &str) -> String {
    static RE_PUNCT: Lazy<Regex> = Lazy::new(|| Regex::new(r"[^\w\s]").unwrap());
    static RE_WHITESPACE: Lazy<Regex> = Lazy::new(|| Regex::new(r"\s+").unwrap());

    let lower = text.to_lowercase();
    let no_punct = RE_PUNCT.replace_all(&lower, "");
    let clean_text = RE_WHITESPACE.replace_all(&no_punct, " ");

    clean_text.trim().to_string()
}

/// Generate a frequency map for a list of words.
fn generate_frequency_map<'a>(words: &[&'a str]) -> HashMap<&'a str, usize> {
    let mut freq_map = HashMap::new();
    for word in words {
        *freq_map.entry(*word).or_insert(0) += 1;
    }
    freq_map
}

/// Calculate term frequency (TF) for a list of unique words and a frequency map.
fn calculate_tf(uniq: &[&str], fm: &HashMap<&str, usize>, total: usize) -> Vec<f64> {
    // Compute TF using the frequency map
    uniq.iter()
        .map(|word| *fm.get(word).unwrap_or(&0) as f64 / total as f64)
        .collect()
}

/// Calculate inverse document frequency (IDF) for a list of unique words and two frequency maps.
fn calculate_idf(uniq: &[&str], fm1: &HashMap<&str, usize>, fm2: &HashMap<&str, usize>) -> Vec<f64> {
    let mut doc_freq = HashMap::new();

    // Populate document frequencies
    for &word in uniq {
        doc_freq.insert(word, fm1.contains_key(word) as usize + fm2.contains_key(word) as usize);
    }

    uniq.iter()
        .map(|word| {
            let count = *doc_freq.get(word).unwrap_or(&0);
            (1.0 + 2.0 / (count as f64 + 1.0)).ln()
        })
        .collect()
}

/// Calculate the TF-IDF for a list of term frequencies and a list of inverse document frequencies.
fn calculate_tf_idf(tf: &[f64], idf: &[f64]) -> Vec<f64> {
    tf.iter().zip(idf.iter()).map(|(a, b)| a * b).collect()
}

/// Calculate the cosine similarity between two vectors.
fn calculate_similarity(tf_idf1: &[f64], tf_idf2: &[f64]) -> f64 {
    let dot_product: f64 = tf_idf1.iter().zip(tf_idf2.iter()).map(|(a, b)| a * b).sum();

    let norm1 = tf_idf1.iter().map(|a| a * a).sum::<f64>().sqrt();
    let norm2 = tf_idf2.iter().map(|a| a * a).sum::<f64>().sqrt();

    if norm1.abs() < f64::EPSILON || norm2.abs() < f64::EPSILON {
        return 0.0;
    }

    dot_product / (norm1 * norm2)
}

/// Interpret the similarity value as a human-readable string.
fn interpret_similarity(s: f64) -> String {
    match s {
        0.0..=0.2 => "Dissimilar".to_string(),
        0.2..=0.4 => "Slightly Similar".to_string(),
        0.4..=0.6 => "Moderately Similar".to_string(),
        0.6..=0.8 => "Quite Similar".to_string(),
        0.8..=1.0 => "Highly Similar".to_string(),
        _ => "Unknown".to_string(), // Catch-all for unexpected values
    }
}