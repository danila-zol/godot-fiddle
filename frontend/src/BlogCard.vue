<script setup>
defineProps(['post'])
</script>

<template>
    <div class="blog-card">
        <a :href='"/" + post.title'>
            <img class="blog-thumbnail" :src="post.thumb">
            <p id="blog-title">{{ post.title }}</p>
            <p id="blog-description">{{ post.text }}</p>
        </a>
    </div>
</template>

<style lang="css" scoped>
.blog-card {
    display: flex;
    flex-direction: column;
}

.blog-card a {
    text-decoration: none;
    color: black;
    border-style: solid;
}

.blog-thumbnail {
    max-height: 300px;
    max-width: 300px;
}
</style>