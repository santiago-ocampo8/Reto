<template>
  <v-simple-table>
    <template v-slot:default>
      <thead>
        <tr>
          <th class="text-left">Id</th>
          <th class="text-left">Name</th>
          <th class="text-left">Age</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td>{{ buyer.id }}</td>
          <td>{{ buyer.name }}</td>
          <td>{{ buyer.age }}</td>
        </tr>
      </tbody>
      <br />
      <h2>Shopping history</h2>
      <br />
      <thead>
        <tr>
          <th class="text-left">Device</th>
          <th class="text-left">Ip</th>
          <th class="text-left">Products</th>
          <th class="text-left">Buyers with the same IP</th>
        </tr>
      </thead>

      <tbody>
        <tr v-for="item in Buys" :key="item.name">
          <td>{{ item.buyerc[0].device }}</td>
          <td>{{ item.buyerc[0].ip }}</td>
          <td>
            <tbody>
              <tr v-for="j in item.buyerc[0].product" :key="j.name">
                <td>{{ j.price }}</td>
                <td>-</td>
                <td>{{ j.nameproduct }}</td>
              </tr>
            </tbody>
          </td>
          <td>
            <v-btn
              block
              color="red darken-4"
              dark
              @click="getBuyerIP(item.buyerc[0].ip)"
              >Buyers with the same IP</v-btn
            >
          </td>
        </tr>
      </tbody>
      <br />

      <h2>Recommendations</h2>
      <br />
      <thead>
        <tr>
          <th class="text-left">Name</th>
          <th class="text-left">Price</th>
        </tr>
      </thead>

      <tbody>
        <tr v-for="item in recommendations" :key="item.name">
          <td>{{ item.nameproduct }}</td>
          <td>{{ item.price }}</td>
        </tr>
      </tbody>
    </template>
  </v-simple-table>
</template>

<script>
export default {
  name: "detailsBuyer",
  data: () => ({
    buyer: {},
    Buys: [],
    recommendations: [],
  }),
  mounted() {
    this.getBuyer();
  },
  methods: {
    getBuyer() {
      let idAlumno = this.$route.params.id;
      const url = "http://localhost:3000/" + "buyerproducts/" + idAlumno;
      this.axios.get(url).then((response) => {
        var temporal = JSON.stringify(response.data.find_buyer);
        var temporal2 = temporal.replace(/~buyerc/g, "buyerc");

        var final = JSON.parse(temporal2);

        this.buyer = final[0];
        this.Buys = final;
        this.getRecommendations(
          this.Buys[Math.floor(Math.random() * (this.Buys.length - 0) + 0)]
            .buyerc[0].product[0].idproduct
        );
      });
    },
    getBuyerIP(ip) {
      this.$router.push({name:'ListBuyersIp',params:{ip:ip}})
    },
    getRecommendations(idproduct) {
      const url = "http://localhost:3000/" + "product/" + idproduct;
      this.axios.get(url).then((response) => {
        var temporal = JSON.stringify(response.data.find_products);
        var temporal2 = temporal.replace(/~product/g, "products");
        var final = JSON.parse(temporal2);

        for (var i = 0; i < final.length; i++) {
          for (var j = 0; j < final[i].products[0].product.length; j++) {
            this.recommendations.push(final[i].products[0].product[j]);
          }
        }
      });
    },
  },
};
</script>
<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
</style>
