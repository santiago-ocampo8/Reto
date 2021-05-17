<template>
  <v-simple-table>
    <template v-slot:default>
      <thead>
        <tr>
          <th class="text-left">Name</th>
          <th class="text-left">Age</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="item in buyers" :key="item.name">
          <td>{{item.buyerc.name}}</td>
          <td>{{item.buyerc.age}}</td>
        </tr>
      </tbody>
    </template>
  </v-simple-table>
</template>

<script>
export default {
  name: "ipbuyers",
  data: () => ({
    buyers: []
  }),
  mounted() {
    this.getBuyer();
  },
  methods: {
    getBuyer() {
      let ip = this.$route.params.ip;
        const url = "http://localhost:3000/" + "buyersip/" + ip;
      this.axios.get(url).then((response) => {
        this.buyers=response.data.find_ip
      });
    }
   }
   
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
</style>
