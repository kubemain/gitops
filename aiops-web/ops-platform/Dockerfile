FROM node:20-alpine

WORKDIR /app

# 安装依赖
COPY package*.json ./
RUN npm install

# 复制源码
COPY . .

# 构建
RUN npm run build

# 暴露端口
EXPOSE 5173

# 运行预览服务器
CMD ["npm", "run", "preview", "--", "--host", "0.0.0.0", "--port", "5173"]