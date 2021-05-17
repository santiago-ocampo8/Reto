<template>
  <v-app>
    <v-app-bar
      app
      color="red darken-4"
      dark
      
    >
      <div class="d-flex align-center">
        <v-btn
        target="_blank"
        text
        @click="goHome"
      >
        <span class="mr-2"><h2>Report Shop</h2></span>
      </v-btn>
        
      </div>

      <v-spacer></v-spacer>
      <v-btn
        target="_blank"
        text
        @click="loadBuyers"
      >
        <span class="mr-2">Load Info</span>
      </v-btn>
      <v-btn
        target="_blank"
        text
        @click="openBuyers"
      >
        <span class="mr-2">Buyers</span>
      </v-btn>
    </v-app-bar>

    <v-main>
      <v-container>
      <router-view/>
      </v-container>
    </v-main>
  </v-app>
</template>

<script>

export default {
  name: 'App',

  data: () => ({
    //
  }),
  methods:{
    openBuyers(){
      this.$router.push({name:'Buyers'})
    },
    loadBuyers(){
      this.$router.push({name:'LoadData'})
    },
    goHome(){
        this.$router.push({name:'Home'})
    }
  }
};
</script>
