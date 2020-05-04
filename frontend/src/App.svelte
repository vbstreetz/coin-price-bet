<script>
  import {Router} from "@sveltech/routify";
  import {routes} from "@sveltech/routify/tmp/routes";
</script>

<style>
  .topbar {
    padding: 20px 0;
  }
  a {
    color: #888;
  }
  a:visited,
  a:hover {
    opacity: 0.6;
  }
</style>

<div class="topbar flex justify-end text-xs">
  <a href="/" class="mr-3">bet</a>
  <a href="/gold" class="mr-3">goldcdp</a>
  <a href="/charts">prices</a>
</div>
<Router {routes}/>
