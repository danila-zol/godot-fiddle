<script setup>
import { blogs } from './globals';
import BlogCard from './BlogCard.vue';
</script>

<template>
<div id="blogs">
    <BlogCard v-for="blog in blogs" :post="blog"></BlogCard>
</div>
</template>

<style lang="css" scoped>
#blogs {
    display: grid;
    grid-template-columns: 1fr 1fr 1fr;
}
</style>