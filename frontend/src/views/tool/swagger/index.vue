<template>
  <div class="page-stack swagger-page">
    <a-card class="panel-card swagger-card" :bordered="false">
      <iframe
        class="swagger-frame"
        title="Swagger UI"
        src="http://192.168.31.81:11080/swagger/"
      />
    </a-card>
  </div>
</template>

<style scoped>
.swagger-page {
  min-height: calc(100vh - 120px);
  display: flex;
}

.swagger-card {
  flex: 1;
  display: flex;
  min-height: 0;
}

:deep(.swagger-card .arco-card-body) {
  flex: 1;
  display: flex;
  padding: 0;
  min-height: 0;
}

.swagger-frame {
  flex: 1;
  width: 100%;
  min-height: calc(100vh - 160px);
  border: 0;
  display: block;
}
</style>
