<script setup>
import { ref, watch } from 'vue'
import Navbar from './Navbar.vue';
import SectionAdder from './SectionAdder.vue';
import Blogpost from './Blogposts.vue';

let pass = import.meta.env.VITE_PASSWORD
</script>

<template>
    <Navbar />
    <SectionAdder />
    <Blogpost />
</template>